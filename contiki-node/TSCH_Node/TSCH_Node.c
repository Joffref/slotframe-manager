#include "contiki.h"
#include "net/netstack.h"
#include "net/nullnet/nullnet.h"
#include "net/ipv6/simple-udp.h"
#include "net/ipv6/uip-ds6-nbr.h"
#include "net/linkaddr.h"
#include <stdio.h>
#include <string.h>

// Define the UDP port to use
#define UDP_PORT 12345

// Process to handle sending data
PROCESS(udp_send_process, "UDP Send Process");
AUTOSTART_PROCESSES(&udp_send_process);

// UDP connection
static struct simple_udp_connection udp_connection;

// Declare the et timer
static struct etimer et;

PROCESS_THREAD(udp_send_process, ev, data)
{
  PROCESS_BEGIN();

  // Initialize nullnet
  nullnet_set_input_callback(NULL);

  // Initialize UDP connection
  simple_udp_register(&udp_connection, UDP_PORT, NULL, UDP_PORT, NULL);

  // Main loop
  while (1)
  {
    // Your data to send
    // For simplicity, we assume that ID is a string.
    const char *id = "123";

    // MAC address
    linkaddr_t mac_address;
    linkaddr_copy(&mac_address, &linkaddr_node_addr);

    // Parent MAC address
    linkaddr_t parent_lladdr;

    // Iterate through the neighbor list to find the parent
    uip_ds6_nbr_t *nbr;
    for (nbr = uip_ds6_nbr_head(); nbr != NULL; nbr = uip_ds6_nbr_next(nbr))
    {
      if (nbr->state == NBR_REACHABLE)
      {
        // Found a reachable neighbor, use its lladdr
        const uip_lladdr_t *lladdr = uip_ds6_nbr_get_ll(nbr);
        // Convert uip_lladdr_t to linkaddr_t
        linkaddr_copy(&parent_lladdr, (const linkaddr_t *)lladdr);
        break;
      }
    }

    // Check if a parent lladdr was found
    if (nbr != NULL)
    {
      // Combine the data into a single message
      char message[100];
      snprintf(message, sizeof(message), "MAC: %d.%d, Parent MAC: %d.%d, ID: %s",
               mac_address.u8[0], mac_address.u8[1], parent_lladdr.u8[0], parent_lladdr.u8[1], id);

      // Send the data
      simple_udp_sendto(&udp_connection, message, strlen(message) + 1, uip_ds6_nbr_get_ipaddr(nbr));
    }

    // Wait for a while before sending again
    etimer_set(&et, CLOCK_SECOND * 10);
    PROCESS_WAIT_EVENT_UNTIL(etimer_expired(&et));
  }

  PROCESS_END();
}
