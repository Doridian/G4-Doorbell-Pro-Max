
#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <unistd.h>

#include "libbmkt/bmkt.h"
#include "libbmkt/custom.h"

#define BMKT_WRAP(FUNC) { \
        int res = FUNC; \
        if (res != BMKT_SUCCESS) { \
            printf(#FUNC " failed (%d)\n", res); \
            exit_program(session, 1); \
        } else { \
            printf(#FUNC " OK\n"); \
        } \
    }

void exit_program(bmkt_ctx_t* session, int code) {
    if (session != NULL) {
        BMKT_WRAP(bmkt_close(session));
        BMKT_WRAP(bmkt_exit(session));
    }
    exit(code);
}

void run_bmkt_identify(bmkt_ctx_t* session) {
    BMKT_WRAP(bmkt_identify(session));
}

void run_bmkt_enroll(bmkt_ctx_t* session, const char* user_id, int finger_id) {
    BMKT_WRAP(bmkt_enroll(session, user_id, strlen(user_id), finger_id));
}

void on_enroll_progress(int percent) {
    printf("Enrollment progress %d %%\n", percent);
}

int on_response(bmkt_response_t* resp, void* cb_ctx) {
    switch (resp->response_id) {
        // Events
        case BMKT_EVT_FINGER_REPORT:
            switch (resp->response.finger_event_resp.finger_state) {
                case BMKT_EVT_FINGER_STATE_NOT_ON_SENSOR:
                    printf("Finger removed from sensor!\n");
                    break;
                case BMKT_EVT_FINGER_STATE_ON_SENSOR:
                    printf("Finger placed on sensor, running identify!\n");
                    run_bmkt_identify((bmkt_ctx_t*)cb_ctx);
                    break;
            }
            break;

        // Enrollment
        case BMKT_RSP_ENROLL_READY:
            on_enroll_progress(0);
            break;
        case BMKT_RSP_ENROLL_FAIL:
            on_enroll_progress(-1);
            break;
        case BMKT_RSP_ENROLL_OK:
            on_enroll_progress(100);
            break;
        case BMKT_RSP_ENROLL_REPORT:
            on_enroll_progress(resp->response.enroll_resp.progress);
            break;

        // Verify / verify_resp
        case BMKT_RSP_VERIFY_READY:
            printf("Verify started!\n");
            break;
        case BMKT_RSP_VERIFY_OK:
            resp->response.verify_resp.user_id[BMKT_MAX_USER_ID_LEN - 1] = 0; // Just to be safe...
            printf("Verify OK! You are %s finger %d\n", resp->response.verify_resp.user_id, resp->response.verify_resp.finger_id);
            break;
        case BMKT_RSP_VERIFY_FAIL:
            printf("Verify FAIL!\n");
            break;

        // Identify / id_resp
        case BMKT_RSP_ID_READY:
            printf("Identify started!\n");
            break;
        case BMKT_RSP_ID_OK:
            resp->response.id_resp.user_id[BMKT_MAX_USER_ID_LEN - 1] = 0; // Just to be safe...
            printf("Identify OK! You are %s finger %d\n", resp->response.id_resp.user_id, resp->response.id_resp.finger_id);
            break;
        case BMKT_RSP_ID_FAIL:
            printf("Identify FAIL!\n");
            break;

        // Unhandled
        default:
            printf("on_response(%d / 0x%02x)\n", resp->response_id, resp->response_id);
            break;
    }
    return BMKT_SUCCESS;
}

int on_error(uint16_t error, void *cb_ctx) {
    printf("on_error(%d)\n", error);
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
    BMKT_WRAP(bmkt_open(session, &sensor, &session, &on_response, session, &on_error, session));

    sleep(1);
    int bmkt_init_fps_ret;
    while ((bmkt_init_fps_ret = bmkt_init_fps(session)) == BMKT_SENSOR_NOT_READY) {
        usleep(10);
    }
    BMKT_WRAP(bmkt_init_fps_ret);
    // TODO: Wait for init_fps ok response maybe
    printf("BMKT initialized!\n");

    //BMKT_WRAP(bmkt_identify(session));
    //BMKT_WRAP(bmkt_enroll(session, "doridian\0", strlen("doridian"), 1));
    while (1) {
        sleep(1);
    }

    exit_program(session, 0);
    return 0;
}

