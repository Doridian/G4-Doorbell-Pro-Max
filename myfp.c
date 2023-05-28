
#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <unistd.h>

#include "libbmkt/bmkt.h"
#include "libbmkt/custom.h"

#define BMKT_WRAP(FUNC) { \
        int res = FUNC; \
        if (res != BMKT_SUCCESS) { \
            printf(#FUNC " failed (%d)\n", res); \
            exit(1); \
        } else { \
            printf(#FUNC " OK\n"); \
        } \
    }

int on_response(bmkt_response_t *resp, void *cb_ctx) {
    printf("on_response(%d)\n", resp->response_id);
    return BMKT_SUCCESS;
}

int on_error(uint16_t error, void *cb_ctx) {
    printf("on_error(%d)\n", error);
    return BMKT_SUCCESS;
}

int on_event(bmkt_finger_event_t *event, void *cb_ctx) {
    printf("on_event(%d)\n", event->finger_state);
    return BMKT_SUCCESS;
}

int main() {
    printf("Initializing BMKT...\n");

    bmkt_sensor_t sensor;
    sensor.type = 0;
    sensor.info.mode = 0;
    sensor.info.unk3 = 0;
    sensor.info.unk4 = 0;
    sensor.info.unk5 = 0x45;
    sensor.info.gpio_number = 0;
    sensor.info.speed = 4000000;
    sensor.info.bpw = 8;
    sensor.info.unk6 = 2;
    sensor.info.unk7 = 0;
    sensor.info.addr = 1;
    sensor.info.subaddr = 1;
    sensor.info.unk1 = 0x44;
    sensor.info.unk2 = 1;
    sensor.info.unk8 = 0;

    bmkt_ctx_t* session;
    BMKT_WRAP(bmkt_init(&session));
    BMKT_WRAP(bmkt_open(session, &sensor, &session, &on_response, NULL, &on_error, NULL));

    sleep(1);
    int bmkt_init_fps_ret;
    while ((bmkt_init_fps_ret = bmkt_init_fps(session)) == BMKT_SENSOR_NOT_READY) {
        usleep(10);
    }
    BMKT_WRAP(bmkt_init_fps_ret);
    // TODO: Wait for init_fps ok response
    sleep(1);
    BMKT_WRAP(bmkt_identify(session));
    // TODO: Actually do something lol
    sleep(1);

    printf("BMKT initialized!\n");

    return 0;
}

