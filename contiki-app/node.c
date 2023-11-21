#include "contiki.h"
#include <stdio.h>
#include "app.h"
/*---------------------------------------------------------------------------*/
PROCESS(node_process, "Node process");
AUTOSTART_PROCESSES(&node_process);
/*---------------------------------------------------------------------------*/
PROCESS_THREAD(node_process, ev, data)
{
  PROCESS_BEGIN();

  printf("Go code returns: %u\n", HelloWorld());

  PROCESS_END();
}
/*---------------------------------------------------------------------------*/