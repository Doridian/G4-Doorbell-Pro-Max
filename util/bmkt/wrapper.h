#pragma once

#include <stdint.h>
#include <libbmkt/custom.h>

#ifndef NULL
#define NULL 0
#endif

extern void go_bmkt_on_error(uint64_t id, uint16_t code);
extern void go_bmkt_on_response(uint64_t id, bmkt_response_t* resp);

static int c_bmkt_on_response(bmkt_response_t* resp, void* cb_ctx_void) {
    uint64_t id = *(uint64_t*)cb_ctx_void;
    go_bmkt_on_response(id, resp);
    return BMKT_SUCCESS;
}

static int c_bmkt_on_error(uint16_t code, void* cb_ctx_void) {
    uint64_t id = *(uint64_t*)cb_ctx_void;
    go_bmkt_on_error(id, code);
    return BMKT_SUCCESS;
}

static bmkt_ctx_t* bmkt_wrapped_init() {
    bmkt_ctx_t* session;
    if (bmkt_init(&session) != BMKT_SUCCESS) {
        return NULL;
    }
    return session;
}

static int bmkt_wrapped_open(bmkt_ctx_t* session, bmkt_sensor_t* sensor, uint64_t* id) {
    bmkt_ctx_t* session_out;
    return bmkt_open(session, sensor, &session_out, &c_bmkt_on_response, id, &c_bmkt_on_error, id);
}
