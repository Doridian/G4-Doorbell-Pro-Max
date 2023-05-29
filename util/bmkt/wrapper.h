#pragma once

#include <stdlib.h>
#include <unistd.h>

#include <libbmkt/bmkt.h>
#include <libbmkt/custom.h>

extern void c_on_error(uint64_t id, uint16_t code);
extern void c_on_response(uint64_t id, bmkt_response_t* resp);

static int on_response(bmkt_response_t* resp, void* cb_ctx_void) {
    uint64_t id = *(uint64_t*)cb_ctx_void;
    c_on_response(id, resp);
    return BMKT_SUCCESS;
}

static int on_error(uint16_t code, void* cb_ctx_void) {
    uint64_t id = *(uint64_t*)cb_ctx_void;
    c_on_error(id, code);
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
    return bmkt_open(session, sensor, &session_out, &on_response, id, &on_error, id);
}
