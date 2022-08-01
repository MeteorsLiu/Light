#ifndef _LIGHT_H
#define _LIGHT_H

extern void UPLOAD(uintptr_t rule, char * ip, char * rates);
//Global Var
typedef struct LinkedList * List;
struct LinkedList
{
	uint32_t Key;
	uint32_t usedTime;
	bpf_u_int32 bytes;
	time_t record;
	List   prev;
	List   next;
};
static pcap_t *p = NULL;
static List Head = NULL;
static List Last = NULL;
static uint32_t maxTimes = 0;
static time_t lastRecycleTime = 0;
static uintptr_t ruleHandle;
//Macro
#define ERROR -1
#define SECOND 1
#define HOUR   3600
#define MiB    131072
#define KiB	   128

void Create(uint32_t Key, bpf_u_int32 bytes);
void Insert(uint32_t Key, bpf_u_int32 bytes);
void Adjust(List Node);
void RecycleLinkedList();
bpf_u_int32 Update(uint32_t Key, bpf_u_int32 bytes);
void Free();
void stop_capture(int o);
int create_pcap_handle(char* device, char* filter);
void packetHandler(
    u_char *args,
    const struct pcap_pkthdr* header,
    const u_char* packet
);

void Init(uintptr_t rule, const char *devName, const char *rule);
#endif
