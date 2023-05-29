#pragma once

#include <stdint.h>

#include "bmkt.h"
#include "bmkt_response.h"
#include "bmkt_message.h"

#define BMKT_MAX_PENDING_SESSIONS 2

// Parts of these are from synaptics libfprint fork or libfprint itself

#define bmkt_ctx_t void

typedef enum bmkt_sensor_state {
    BMKT_SENSOR_STATE_UNINIT = 0,
    BMKT_SENSOR_STATE_IDLE,
    BMKT_SENSOR_STATE_INIT,
    BMKT_SENSOR_STATE_EXIT,
} bmkt_sensor_state_t;

typedef enum {
    BMKT_OP_STATE_START = -1,
    BMKT_OP_STATE_GET_RESP,
    BMKT_OP_STATE_WAIT_INTERRUPT,
    BMKT_OP_STATE_SEND_ASYNC,
    BMKT_OP_STATE_COMPLETE,
} bmkt_op_state_t;

typedef int (*bmkt_resp_cb_t)(bmkt_response_t *resp, void *cb_ctx);
typedef int (*bmkt_event_cb_t)(bmkt_finger_event_t *event, void *cb_ctx);
typedef int (*bmkt_general_error_cb_t)(uint16_t error, void *cb_ctx);

typedef struct bmkt_sensor_version {
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

typedef struct bmkt_session_ctx {
    uint8_t seq_num;
    bmkt_resp_cb_t resp_cb;
    void *cb_ctx;
} bmkt_session_ctx_t;

typedef enum  {
    GPIO_DIRECTION_IN = 0,
    GPIO_DIRECTION_OUT = 1,
} gpio_direction_t;

typedef enum {
    GPIO_EDGE_NONE = 0,
    GPIO_EDGE_FALLING = 1,
    GPIO_EDGE_RISING = 2,
    GPIO_EDGE_BOTH = 3,
} gpio_edge_t;

typedef struct gpio_transport_info_t {
    int pin;
    gpio_direction_t direction;
    gpio_edge_t edge;
    int active_low;
} gpio_transport_info_t;

#define SPI_CPHA 0x01 /* clock phase */
#define SPI_CPOL 0x02 /* clock polarity */
typedef enum {
    SPI_MODE_0 = (0|0),            /* (original MicroWire) */
    SPI_MODE_1 = (0|SPI_CPHA),
    SPI_MODE_2 = (SPI_CPOL|0),
    SPI_MODE_3 = (SPI_CPOL|SPI_CPHA),
} spi_mode_t;

typedef enum {
    SENSOR_TRANSPORT_SPI = 0
} sensor_transport_t;

typedef struct spi_transport_info {
    spi_mode_t mode;
    int speed;
    int bpw;
    int addr;
    int subaddr;
    gpio_transport_info_t pin_out;
    gpio_transport_info_t pin_in;
    int unknown_padding;
} spi_transport_info_t;

typedef struct bmkt_sensor {
    sensor_transport_t transport_type;
    spi_transport_info_t transport_info;

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
int bmkt_delete_enrolled_user(bmkt_ctx_t* session, uint8_t finger_id, const uint8_t* user_id,  uint32_t user_id_len);
int bmkt_enroll(bmkt_ctx_t* session, const uint8_t* user_id,  uint32_t user_id_len, uint8_t finger_id);

int bmkt_close(bmkt_ctx_t* session);
int bmkt_exit(bmkt_ctx_t* session);
int bmkt_cancel_op(bmkt_ctx_t* session);
