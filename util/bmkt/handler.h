#pragma once

#include <stdlib.h>

#include <libbmkt/bmkt.h>
#include <libbmkt/custom.h>

typedef enum {
    IF_STATE_INVALID = -1,
    IF_STATE_IDLE = 0,
    IF_STATE_INIT,
    IF_STATE_ENROLLING,
    IF_STATE_VERIFYING,
    IF_STATE_IDENTIFYING,
    IF_STATE_CANCELLING,
    IF_STATE_DELETING_ALL,
} cb_ctx_state_t; 

typedef struct cb_ctx_struct {
    bmkt_ctx_t* session;
    cb_ctx_state_t state;
    int last_error;
} cb_ctx_t;

int on_response(bmkt_response_t* resp, void* cb_ctx_void) {
    return BMKT_SUCCESS;
}

int on_error(uint16_t error, void* cb_ctx_void) {
    return BMKT_SUCCESS;
}

static void bmkt_main_close(cb_ctx_t* ctx) {
    if (ctx->session) {
        bmkt_session_ctx_t* session = ctx->session;
        ctx->session = NULL;
        bmkt_close(session);
        bmkt_exit(session);
    }
    free(ctx);
}

static cb_ctx_t* bmkt_main_init() {
    bmkt_sensor_t sensor;
    // Type 0 means SPI in this library
    sensor.transport_type = SENSOR_TRANSPORT_SPI;

    // SPI settings
    sensor.transport_info.addr = 1;
    sensor.transport_info.subaddr = 1;
    sensor.transport_info.mode = SPI_MODE_0;
    sensor.transport_info.speed = 4000000;
    sensor.transport_info.bpw = 8;

    // GPIO pin information
    sensor.transport_info.pin_out.pin = 68;
    sensor.transport_info.pin_out.direction = GPIO_DIRECTION_OUT;
    sensor.transport_info.pin_out.edge = GPIO_EDGE_NONE;
    sensor.transport_info.pin_out.active_low = 0;

    sensor.transport_info.pin_in.pin = 69;
    sensor.transport_info.pin_in.direction = GPIO_DIRECTION_IN;
    sensor.transport_info.pin_in.edge = GPIO_EDGE_RISING;
    sensor.transport_info.pin_in.active_low = 0;

    // No idea, might just be padding
    sensor.transport_info.unknown_padding = 0;

    bmkt_ctx_t* session;
    bmkt_ctx_t* session_out;
    cb_ctx_t* ctx = malloc(sizeof(cb_ctx_t));
    int res;

    res = bmkt_init(&session);
    if (res != BMKT_SUCCESS) {
        ctx->state = IF_STATE_INVALID;
        ctx->last_error = res;
        return ctx;
    }

    ctx->session = session;
    ctx->state = IF_STATE_INIT;
    res = bmkt_open(session, &sensor, &session_out, &on_response, &ctx, &on_error, &ctx);
    if (res != BMKT_SUCCESS) {
        ctx->state = IF_STATE_INVALID;
        ctx->last_error = res;
        return ctx;
    }

    return ctx;
}
