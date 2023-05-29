
#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

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
    printf("on_response(%d / 0x%02x)\n", resp->response_id, resp->response_id);
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

/*
    GPIO
    [0] ID
    [1] DIRECTION 0=IN 1=OUT
    [2] EDGE 2=rising 3=both 1=falling 0=none
    [3] ACTIVE_LOW 0=0 1=1
*/

int main() {
    printf("Initializing BMKT...\n");

    bmkt_sensor_t sensor;
    // Type 0 means SPI in this library
    sensor.type = 0;

    // SPI settings
    sensor.info.addr = 1;
    sensor.info.subaddr = 1;
    sensor.info.mode = SPI_MODE_0;
    sensor.info.speed = 4000000;
    sensor.info.bpw = 8;

    // GPIO pin information
    sensor.info.pin1.pin = 68;
    sensor.info.pin1.direction = GPIO_DIRECTION_OUT;
    sensor.info.pin1.edge = GPIO_EDGE_NONE;
    sensor.info.pin1.active_low = 0;

    sensor.info.pin2.pin = 69;
    sensor.info.pin2.direction = GPIO_DIRECTION_IN;
    sensor.info.pin2.edge = GPIO_EDGE_RISING;
    sensor.info.pin2.active_low = 0;

    // No idea, might just be padding
    sensor.info.unknown_padding = 0;

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
    //BMKT_WRAP(bmkt_enroll(session, "doridian\0", strlen("doridian"), 1));
    // TODO: Actually do something lol
    while (1) {
        sleep(1);
    }

    printf("BMKT initialized!\n");

    return 0;
}

