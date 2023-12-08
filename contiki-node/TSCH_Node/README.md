# TSCH Node - Sending Data over UDP

> Note this is compiling but not running on the node yet.

Ce code Contiki est destiné à être utilisé sur des nœuds IoT (Internet des Objets) utilisant le protocole IEEE 802.15.4e TSCH (Time-Slotted Channel Hopping). Il envoie des données sur un réseau IPv6 en utilisant le protocole UDP (User Datagram Protocol).

# Fonctionnalités :

    Envoi de données UDP : Le nœud envoie périodiquement des données UDP contenant des informations sur son adresse MAC, l'adresse MAC de son parent et un identifiant unique.

    Utilisation du voisin atteignable : Le code recherche un voisin atteignable dans la liste des voisins pour obtenir l'adresse MAC de son parent.

    Intervalle de transmission configurable : Le nœud attend un certain temps avant d'envoyer à nouveau les données, ce qui est configurable.

# Configuration
Assurez-vous d'avoir Contiki installé sur votre système avant de compiler et d'exécuter ce code.

# Compilation
Utilisez la commande make pour compiler le code.

# Personnalisation
Vous pouvez personnaliser le code en modifiant les paramètres tels que le port UDP, l'intervalle de transmission, etc. Ces paramètres sont définis dans le fichier TSCH_Node.c.
