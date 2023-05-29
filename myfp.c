
#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <unistd.h>

#include "libbmkt/bmkt.h"
#include "libbmkt/custom.h"

#define BMKT_WRAP(FUNC, session) { \
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
        bmkt_close(session);
        bmkt_exit(session);
    }
    exit(code);
}

typedef enum {
    IF_STATE_INIT = 0,
    IF_STATE_IDLE,
    IF_STATE_ENROLLING,
    IF_STATE_VERIFYING,
    IF_STATE_IDENTIFYING,
    IF_STATE_CANCELLING,
} cb_ctx_state_t; 

typedef struct cb_ctx_struct {
    bmkt_ctx_t* session;
    cb_ctx_state_t state;
} cb_ctx_t;

static inline void bmkt_cancel_if_running(cb_ctx_t* ctx) {
    ctx->state = IF_STATE_CANCELLING;
    BMKT_WRAP(bmkt_cancel_op(ctx->session), ctx->session);
    // TODO: Wait for cancellation completion
}

void run_bmkt_identify(cb_ctx_t* ctx) {
    bmkt_cancel_if_running(ctx);
    ctx->state = IF_STATE_IDENTIFYING;
    BMKT_WRAP(bmkt_identify(ctx->session), ctx->session);
}

void run_bmkt_enroll(cb_ctx_t* ctx, const char* user_id, int finger_id) {
    bmkt_cancel_if_running(ctx);
    ctx->state = IF_STATE_ENROLLING;
    BMKT_WRAP(bmkt_enroll(ctx->session, (const uint8_t*)user_id, strlen(user_id), finger_id), ctx->session);
}

void on_enroll_progress(int percent) {
    printf("Enrollment progress %d %%\n", percent);
}

// IDENTIFY/VERIFY/ENROLL CANNOT BE RUN AFTER FINGER IS ALREADY ON SENSOR!
int on_response(bmkt_response_t* resp, void* cb_ctx_void) {
    cb_ctx_t* ctx = (cb_ctx_t*)cb_ctx_void;

    switch (resp->response_id) {
        // Events
        case BMKT_EVT_FINGER_REPORT:
            switch (resp->response.finger_event_resp.finger_state) {
                case BMKT_EVT_FINGER_STATE_NOT_ON_SENSOR:
                    printf("Finger removed from sensor!\n");
                    if (ctx->state == IF_STATE_IDLE) {
                        run_bmkt_identify(ctx);
                    }
                    break;
                case BMKT_EVT_FINGER_STATE_ON_SENSOR:
                    printf("Finger placed on sensor!\n");
                    break;
            }
            break;

        // Init
        case BMKT_RSP_FPS_INIT_FAIL:
            printf("Init FAIL!\nBailing...\n");
            exit(1);
            break;
        case BMKT_RSP_FPS_INIT_OK:
            ctx->state = IF_STATE_IDLE;
            printf("Init OK!\n");
            if (!resp->response.init_resp.finger_presence) {
                run_bmkt_identify(ctx);
            }
            break;

        // Enrollment
        case BMKT_RSP_ENROLL_READY:
            ctx->state = IF_STATE_ENROLLING;
            on_enroll_progress(0);
            break;
        case BMKT_RSP_ENROLL_FAIL:
            ctx->state = IF_STATE_IDLE;
            on_enroll_progress(-1);
            break;
        case BMKT_RSP_ENROLL_OK:
            ctx->state = IF_STATE_IDLE;
            on_enroll_progress(100);
            break;
        case BMKT_RSP_ENROLL_REPORT:
            ctx->state = IF_STATE_ENROLLING;
            on_enroll_progress(resp->response.enroll_resp.progress);
            break;

        // Verify / verify_resp
        case BMKT_RSP_VERIFY_READY:
            ctx->state = IF_STATE_VERIFYING;
            printf("Verify started!\n");
            break;
        case BMKT_RSP_VERIFY_OK:
            ctx->state = IF_STATE_IDLE;
            resp->response.verify_resp.user_id[BMKT_MAX_USER_ID_LEN - 1] = 0; // Just to be safe...
            printf("Verify OK! You are %s finger %d\n", resp->response.verify_resp.user_id, resp->response.verify_resp.finger_id);
            break;
        case BMKT_RSP_VERIFY_FAIL:
            ctx->state = IF_STATE_IDLE;
            printf("Verify FAIL!\n");
            break;

        // Identify / id_resp
        case BMKT_RSP_ID_READY:
            ctx->state = IF_STATE_IDENTIFYING;
            printf("Identify started!\n");
            break;
        case BMKT_RSP_ID_OK:
            ctx->state = IF_STATE_IDLE;
            resp->response.id_resp.user_id[BMKT_MAX_USER_ID_LEN - 1] = 0; // Just to be safe...
            printf("Identify OK! You are %s finger %d\n", resp->response.id_resp.user_id, resp->response.id_resp.finger_id);
            break;
        case BMKT_RSP_ID_FAIL:
            ctx->state = IF_STATE_IDLE;
            printf("Identify FAIL!\n");
            break;

        case BMKT_RSP_CANCEL_OP_OK:
            ctx->state = IF_STATE_IDLE;
            printf("Canellation OK!\n");
            break;
        case BMKT_RSP_CANCEL_OP_FAIL:
            printf("Cancellation FAIL!\nBailing...\n");
            exit(1);
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

    BMKT_WRAP(bmkt_init(&session), session);
    cb_ctx_t ctx;
    ctx.session = session;
    ctx.state = IF_STATE_INIT;
    BMKT_WRAP(bmkt_open(session, &sensor, &session, &on_response, &ctx, &on_error, &ctx), session);

    sleep(1);
    int bmkt_init_fps_ret;
    while ((bmkt_init_fps_ret = bmkt_init_fps(session)) == BMKT_SENSOR_NOT_READY) {
        usleep(10);
    }
    BMKT_WRAP(bmkt_init_fps_ret, session);
    printf("BMKT initialized!\n");

    while (1) {
        sleep(1);
    }

    exit_program(session, 0);
    return 0;
}

