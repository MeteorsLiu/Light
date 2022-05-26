#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <pcap/pcap.h>
#include <signal.h>
#include <linux/if_ether.h>
#include <netinet/ip.h>
#include <netinet/in.h>
#include <netinet/if_ether.h>
#include <arpa/inet.h> 
#include <time.h>

#include <light.h>


void Create(uint32_t Key, bpf_u_int32 bytes) {
	Head = Last = (List) malloc(sizeof(struct LinkedList));
	Head->Key = Key;
	Head->bytes = bytes;
	Head->record = time(NULL);
	Head->prev = Head->next = NULL;
	Head->usedTime = 0;
}
void Insert(uint32_t Key, bpf_u_int32 bytes) {
	//impossible
	if (Last == NULL)
		stop_capture(1);
	List indirect = NULL;
	indirect = Last->next = (List) malloc(sizeof(struct LinkedList));
	indirect->Key = Key;
	indirect->bytes = bytes;
	indirect->record = time(NULL);
	indirect->prev = Last;
	indirect->next = NULL;
	indirect->usedTime = 0;
	Last = indirect;
}


void Adjust(List Node) {
	maxTimes = Node->usedTime;
	//printf("%d Adjust to be Top!\n", Node->Key);
	//Make it be first
	//First connected the previous node with the next node
	if (Last == Node)
		Last = Node->prev;
	Node->prev->next = Node->next;
	//Update the new next ptr
	Node->next = Head;
	//Update the head prev ptr
	Head->prev = Node;
	Node->prev = NULL;
	Head = Node;
	Node->usedTime = 0;
}

void RecycleLinkedList() {
	time_t nowTime = time(NULL);
	if (nowTime - lastRecycleTime < 4*HOUR)
		return;
	//Don't recycle the head ptr
	List indirect = Head->next;
	List tmp = NULL;
	while (indirect) {
		//Free the inactive resource
		if (nowTime - indirect->record >= 3*HOUR) {
			if (indirect->prev)
				indirect->prev->next = indirect->next;
			if (Last == indirect)
				Last = indirect->prev;
			tmp = indirect;
			indirect = indirect->next;
			free(tmp);
		} else {
			indirect = indirect->next;
		}	
	}
	lastRecycleTime = nowTime;
}
bpf_u_int32 Update(uint32_t Key, bpf_u_int32 bytes) {
	List indirect = Head;
	bpf_u_int32 rate = 0;
	int flag = 0;
	time_t record;
	time_t diff;
	while (indirect) {
		//Matched !
		if (indirect->Key == Key) {
			flag = 1;
			record = time(NULL);
			diff = record - indirect->record;
			if (diff == SECOND) {
				//Reset Counter
				indirect->record = record;
				rate = indirect->bytes;
				indirect->bytes = 0;
				indirect->usedTime++;
				//Compare the maxTimes
				//And Head Node will not be adjusted
				if (indirect->usedTime > maxTimes 
					&& indirect->prev != NULL 
					&& Head != indirect) {
					Adjust(indirect);
					//printf("%d is the Head Node\n", Head->Key);
				}
				return rate;
			} else if (diff > SECOND || diff < 0) {
				//invalid/impossible
				//so reset the counter
				indirect->record = record;
				indirect->bytes = bytes;
			} else {
				indirect->bytes += bytes;
			}
			break;
		}
		indirect = indirect->next;
	}
	if (!flag)
		Insert(Key, bytes);

	return 0;
}

void Free() {
	List indirect = Head;
	List tmp = NULL;
	while (indirect) {
		tmp = indirect->next;
		//printf("Free %d\n", indirect->Key);
		free(indirect);
		indirect = tmp;
	}
	Head = Last = NULL;
	lastRecycleTime = maxTimes = 0;
}
void stop_capture(int o) {
	printf("Exit\n");
	pcap_close(p);
	p = NULL;
	Free();
	printf("Resource clear\n");
	exit(1);
}

int create_pcap_handle(char* device, char* filter)
{
    char errbuf[PCAP_ERRBUF_SIZE];
    struct bpf_program bpf;
    bpf_u_int32 netmask;
    bpf_u_int32 srcip;

    // If no network interface (device) is specfied, get the first one.
    if (device == NULL) {
    	return ERROR;
    }

    // Get network device source IP address and netmask.
    if (pcap_lookupnet(device, &srcip, &netmask, errbuf) == PCAP_ERROR) {
        fprintf(stderr, "pcap_lookupnet: %s\n", errbuf);
        return ERROR;
    }

    // Open the device for live capture.
    p = pcap_open_live(device, BUFSIZ, 0, 1000, errbuf);
    if (p == NULL) {
        fprintf(stderr, "pcap_open_live(): %s\n", errbuf);
        return ERROR;
    }

    // Convert the packet filter epxression into a packet filter binary.
    if (pcap_compile(p, &bpf, filter, 0, netmask) == PCAP_ERROR) {
        fprintf(stderr, "pcap_compile(): %s\n", pcap_geterr(p));
        return ERROR;
    }

    // Bind the packet filter to the libpcap p.
    if (pcap_setfilter(p, &bpf) == PCAP_ERROR) {
        fprintf(stderr, "pcap_setfilter(): %s\n", pcap_geterr(p));
        return ERROR;
    }

    //All done :)
    return 0;
}


void packetHandler(
    u_char *args,
    const struct pcap_pkthdr* header,
    const u_char* packet
) {
    struct ether_header *eth_header;
    bpf_u_int32 rate = 0;
    eth_header = (struct ether_header *) packet;
    //Only IP Packet
    if (ntohs(eth_header->ether_type) == ETHERTYPE_IP) {
    	//Get IP Header
    	//We only need the src ip
    	//no need to parse the tcp header in advance
    	struct ip * ip_header = (struct ip *) (packet + sizeof(struct ether_header));
    	//Create Linked List
    	if (!Head) {
    		Create((ip_header->ip_src).s_addr, header->len);
    		return;
    	}
    	rate = Update((ip_header->ip_src).s_addr, header->len);
    	if (rate > 0) {
    		if (rate >= MiB) {
    			rate /= MiB;
    			printf("record:  %s  %luMbps\n", inet_ntoa(ip_header->ip_src), rate);
    		} else {
    			rate /= KiB;
    			printf("record:  %s  %luKbps\n", inet_ntoa(ip_header->ip_src), rate);
    		}
    		RecycleLinkedList(); 
    	}
    	//printf("ip len : %d\n", ip_header->ip_len);
    }
}

void Init(const char *devName, const char *rule) {
	if (create_pcap_handle(devName, rule) == ERROR)
    	exit(0);
    
    // Initialize the first recycle time.
    lastRecycleTime = time(NULL);
    // Start Loop
    pcap_loop(p, 0, packetHandler, NULL);
}