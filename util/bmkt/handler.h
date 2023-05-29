#pragma once

#include <stdlib.h>
#include <unistd.h>

#include <libbmkt/bmkt.h>
#include <libbmkt/custom.h>

static int on_response(bmkt_response_t* resp, void* cb_ctx_void) {
    return BMKT_SUCCESS;
}

static int on_error(uint16_t error, void* cb_ctx_void) {
    return BMKT_SUCCESS;
}

static bmkt_ctx_t* bmkt_wrapped_init() {
    bmkt_ctx_t* session;
    if (bmkt_init(&session) != BMKT_SUCCESS) {
        return NULL;
    }
    return session;
}

static int bmkt_wrapped_open(bmkt_ctx_t* session, bmkt_sensor_t* sensor) {
    bmkt_ctx_t* session_out;
    return bmkt_open(session, sensor, &session_out, &on_response, NULL, &on_error, NULL);
}
