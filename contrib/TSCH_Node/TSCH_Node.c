#include "contiki.h"
#include "node-id.h"
#include "net/routing/rpl-classic/rpl.h"
#include "net/routing/rpl-classic/rpl-dag.h"

#include <stdio.h> 

 struct slot {
 	uint8_t id;
 	uint8_t channel;
 	struct slot* next;
 };
 typedef struct slot slot;

  struct node {
  	slot* emitting_slots; /* emmiting Slots Frames */
	slot* reciving_slots; /* receiving Slots Frames */
	uint8_t node_id; /* My ID */
	linkaddr_t parent_addr; /* My parent ID */
	uint8_t sf_id;	/* Slot Frame ID for synchro with SFmanager */
	};
	typedef struct node node;
	
  static struct etimer timer;

rpl_parent_t *parent;

PROCESS(packet_transmission,"Packet transmission is on");
AUTOSTART_PROCESSES(&packet_transmission);

PROCESS_THREAD(packet_transmission, ev, data)
{


  PROCESS_BEGIN();

  /* Setup a periodic timer that expires after 10 seconds. */
  etimer_set(&timer, CLOCK_SECOND * 10);
  
  while(1) {
    printf("my node id is %d\n", node_id);
    printf("my parent id is %3u\n",rpl_get_parent_lladdr(parent)->u8[15]);
    /* Wait for the periodic timer to expire and then restart the timer. */
    PROCESS_WAIT_EVENT_UNTIL(etimer_expired(&timer));
    etimer_reset(&timer);
 } 
 PROCESS_END();

}

 


