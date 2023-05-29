#pragma once

#include <stdint.h>

#define BMKT_MAX_PENDING_SESSIONS 2

// Parts of these are from synaptics libfprint fork or libfprint itself

#define bmkt_ctx_t void

typedef enum bmkt_sensor_state
{
	BMKT_SENSOR_STATE_UNINIT 			= 0,
	BMKT_SENSOR_STATE_IDLE,
	BMKT_SENSOR_STATE_INIT,
	BMKT_SENSOR_STATE_EXIT,
} bmkt_sensor_state_t;

typedef enum
{
	BMKT_OP_STATE_START    = -1,
	BMKT_OP_STATE_GET_RESP,
	BMKT_OP_STATE_WAIT_INTERRUPT,
	BMKT_OP_STATE_SEND_ASYNC,
	BMKT_OP_STATE_COMPLETE,
} bmkt_op_state_t;

typedef int (*bmkt_resp_cb_t)(bmkt_response_t *resp, void *cb_ctx);
typedef int (*bmkt_event_cb_t)(bmkt_finger_event_t *event, void *cb_ctx);
typedef int (*bmkt_general_error_cb_t)(uint16_t error, void *cb_ctx);

typedef struct bmkt_sensor_version
{
	uint32_t build_time;
	uint32_t build_num;
	uint8_t version_major;
	uint8_t version_minor;
	uint8_t target;
	uint8_t product;
	uint8_t silicon_rev;
	uint8_t formal_release;
	uint8_t platform;
	uint8_t patch;
	uint8_t serial_number[6];
	uint16_t security;
	uint8_t iface;
	uint8_t device_type;
} bmkt_sensor_version_t;

typedef struct bmkt_session_ctx
{
	uint8_t seq_num;
	bmkt_resp_cb_t resp_cb;
	void *cb_ctx;
} bmkt_session_ctx_t;

typedef struct spi_transport_info {
    uint32_t mode;
    uint32_t speed;
    uint32_t bpw;
    uint32_t addr;
    uint32_t subaddr;
    uint32_t unk1;
    uint32_t unk2;
    uint32_t unk3;
    uint32_t unk4;
    uint32_t unk5;
    uint32_t unk6;
    uint32_t unk7;
    uint32_t unk8;
    uint32_t unk9;
} spi_transport_info_t;

typedef struct bmkt_sensor {
    uint32_t type;
    spi_transport_info_t info;

	bmkt_sensor_version_t version;
	bmkt_session_ctx_t pending_sessions[BMKT_MAX_PENDING_SESSIONS];
	int empty_session_idx;
	int flags;
	int seq_num;
	bmkt_sensor_state_t sensor_state;
	bmkt_event_cb_t finger_event_cb;
	void *finger_cb_ctx;
	bmkt_general_error_cb_t gen_err_cb;
	void *gen_err_cb_ctx;
	bmkt_op_state_t op_state;
} bmkt_sensor_t;

int bmkt_init(bmkt_ctx_t** session);
int bmkt_open(bmkt_ctx_t* session, bmkt_sensor_t* sensor, bmkt_ctx_t** session_out, bmkt_resp_cb_t response_cb, void* response_ctx, bmkt_general_error_cb_t error_cb, void* error_ctx);
int bmkt_init_fps(bmkt_ctx_t* session);

int bmkt_identify(bmkt_ctx_t* session);
int bmkt_enroll(bmkt_ctx_t* session, const uint8_t* user_id,  uint32_t user_id_len, uint8_t finger_id);
