#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>

#include "NfcLibrary/inc/Nfc.h"

static unsigned char discoveryTechnologies[] = {
    MODE_POLL | TECH_PASSIVE_NFCA,
    MODE_POLL | TECH_PASSIVE_NFCF,
    MODE_POLL | TECH_PASSIVE_NFCB,
    MODE_POLL | TECH_PASSIVE_15693,
};

#define NFC_WRAP(FUNC) if (FUNC == NFC_ERROR) { \
        printf(#FUNC " failed\n"); \
        exit(1); \
    }

void ndefPullCB(unsigned char *msg, unsigned short len) {
    printf("NDEF data = ");
    for (int i = 0; i < len; i++) {
        printf("%02x ", msg[i]);
    }
    printf("\n");
}

void readAllCards(NxpNci_RfIntf_t rfInterface) {
    do {
        switch(rfInterface.Protocol) {
        case PROT_T1T:
        case PROT_T2T:
        case PROT_T3T:
        case PROT_ISODEP:
            printf("Found NDEF tag, reading...\n");
            // Process NDEF message read
            NxpNci_ProcessReaderMode(rfInterface, READ_NDEF);
            break;

        case PROT_MIFARE:
            break;

        default:
            break;
        }
    } while(rfInterface.MoreTags && NxpNci_ReaderActivateNext(&rfInterface) == NFC_SUCCESS);

    // Wait for card removal
    NxpNci_ProcessReaderMode(rfInterface, PRESENCE_CHECK);
    
}

int main() {
    printf("Initializing NFC...\n");
    NxpNci_RfIntf_t rfInterface;

    RW_NDEF_RegisterPullCallback((void*)(*ndefPullCB));

    NFC_WRAP(NxpNci_Connect());
    NFC_WRAP(NxpNci_ConfigureSettings());
    NFC_WRAP(NxpNci_ConfigureMode(NXPNCI_MODE_RW));
    NFC_WRAP(NxpNci_StartDiscovery(discoveryTechnologies, sizeof(discoveryTechnologies)));

    printf("NFC initialized!\n");

    while (1) {
        NxpNci_WaitForDiscoveryNotification(&rfInterface);
        readAllCards(rfInterface);
        NFC_WRAP(NxpNci_StopDiscovery());
        NFC_WRAP(NxpNci_StartDiscovery(discoveryTechnologies, sizeof(discoveryTechnologies)));
    }

    return 0;
}
