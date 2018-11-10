//-----------------------------------------------------------------------------
/*

Golang wrapper for jackd2.

*/
//-----------------------------------------------------------------------------

package jack

/*
#cgo linux LDFLAGS: -ljack
#cgo windows,386 LDFLAGS: -llibjack
#cgo windows,amd64 LDFLAGS: -llibjack64

#include <jack/jack.h>
#include <jack/midiport.h>

extern int goProcess(unsigned int, void *);
extern int goBufferSize(unsigned int, void *);
extern int goSampleRate(unsigned int, void *);
extern void goPortRegistration(jack_port_id_t, int, void *);
extern void goPortRename(jack_port_id_t, const char *, const char *, void *);
extern void goPortConnect(jack_port_id_t, jack_port_id_t, int, void *);
extern void goClientRegistration(const char *, int, void *);
extern void goFreewheel(int, void *);
extern int goGraphOrder(void *);
extern int goXrun(void *);
extern void goShutdown(void *);
extern void goInfoShutdown(jack_status_t, const char *, void *);
extern void goErrorFunction(const char *);
extern void goInfoFunction(const char *);

jack_client_t * jack_client_open_go(const char *client_name, int options, int *status) {
  return jack_client_open(client_name, (jack_options_t)options, (jack_status_t *)status);
}

int jack_set_process_callback_go(jack_client_t * client) {
  return jack_set_process_callback(client, goProcess, client);
}

int jack_set_buffer_size_callback_go(jack_client_t * client) {
  return jack_set_buffer_size_callback(client, goBufferSize, client);
}

int jack_set_sample_rate_callback_go(jack_client_t * client) {
	return jack_set_sample_rate_callback(client, goSampleRate, client);
}

int jack_set_port_registration_callback_go(jack_client_t * client) {
	return jack_set_port_registration_callback(client, goPortRegistration, client);
}

int jack_set_port_rename_callback_go(jack_client_t * client) {
	return jack_set_port_rename_callback(client, goPortRename, client);
}

int jack_set_port_connect_callback_go(jack_client_t * client) {
	return jack_set_port_connect_callback(client, goPortConnect, client);
}

int jack_set_client_registration_callback_go(jack_client_t * client) {
  jack_set_client_registration_callback(client, goClientRegistration, client);
}

int jack_set_freewheel_callback_go(jack_client_t * client) {
  jack_set_freewheel_callback(client, goFreewheel, client);
}

int jack_set_graph_order_callback_go(jack_client_t * client) {
  jack_set_graph_order_callback(client, goGraphOrder, client);
}

int jack_set_xrun_callback_go(jack_client_t * client) {
  jack_set_xrun_callback(client, goXrun, client);
}

void jack_on_shutdown_go(jack_client_t * client) {
	jack_on_shutdown(client, goShutdown, client);
}

void jack_on_info_shutdown_go(jack_client_t * client) {
	jack_on_info_shutdown(client, goInfoShutdown, client);
}

void jack_set_error_function_go() {
	jack_set_error_function(goErrorFunction);
}

void jack_set_info_function_go() {
	jack_set_info_function(goInfoFunction);
}

*/
import "C"
import (
	"fmt"
	"strings"
	"sync"
	"unsafe"
)

//-----------------------------------------------------------------------------

// JACK options.
const (
	NullOption    = C.JackNullOption
	NoStartServer = C.JackNoStartServer
	UseExactName  = C.JackUseExactName
	ServerName    = C.JackServerName
	LoadName      = C.JackLoadName
	LoadInit      = C.JackLoadInit
	SessionID     = C.JackSessionID
)

// JACK port flags.
const (
	PortIsInput    = C.JackPortIsInput
	PortIsOutput   = C.JackPortIsOutput
	PortIsPhysical = C.JackPortIsPhysical
	PortCanMonitor = C.JackPortCanMonitor
	PortIsTerminal = C.JackPortIsTerminal
)

// JACK default audio/midi types.
const (
	DefaultAudio = "32 bit float mono audio"
	DefaultMIDI  = "8 bit raw midi"
)

//-----------------------------------------------------------------------------
// status bitfield

// Status is the go type for the JACK status bitmap.
type Status uint

// Status bits.
const (
	Failure       Status = C.JackFailure
	InvalidOption        = C.JackInvalidOption
	NameNotUnique        = C.JackNameNotUnique
	ServerStarted        = C.JackServerStarted
	ServerFailed         = C.JackServerFailed
	ServerError          = C.JackServerError
	NoSuchClient         = C.JackNoSuchClient
	LoadFailure          = C.JackLoadFailure
	InitFailure          = C.JackInitFailure
	ShmFailure           = C.JackShmFailure
	VersionError         = C.JackVersionError
	BackendError         = C.JackBackendError
	ClientZombie         = C.JackClientZombie
)

func (status Status) String() string {
	if status == 0 {
		return ""
	}
	var s []string
	// decode the bits
	statusString := []struct {
		val Status
		str string
	}{
		{Failure, "Failure"},
		{InvalidOption, "InvalidOption"},
		{NameNotUnique, "NameNotUnique"},
		{ServerStarted, "ServerStarted"},
		{ServerFailed, "ServerFailed"},
		{ServerError, "ServerError"},
		{NoSuchClient, "NoSuchClient"},
		{LoadFailure, "LoadFailure"},
		{InitFailure, "InitFailure"},
		{ShmFailure, "ShmFailure"},
		{VersionError, "VersionError"},
		{BackendError, "BackendError"},
		{ClientZombie, "ClientZombie"},
	}
	for _, v := range statusString {
		if status&v.val == v.val {
			s = append(s, v.str)
			status &= ^v.val
		}
	}
	// any leftover bits
	if status != 0 {
		s = append(s, fmt.Sprintf("Unknown(%x)", uint(status)))
	}
	return strings.Join(s, ",")
}

//-----------------------------------------------------------------------------

// GetVersion returns the version of the JACK, in the form of several numbers.
// NOTE returns 0,0,0,0 in 1.9.12
func GetVersion() (int, int, int, int) {
	var major, minor, micro, protocol C.int
	C.jack_get_version(&major, &minor, &micro, &protocol)
	return int(major), int(minor), int(micro), int(protocol)
}

// GetVersionString returns the version of the JACK, in the form of a string.
func GetVersionString() string {
	return C.GoString(C.jack_get_version_string())
}

// ClientNameSize returns the maximum number of characters in a JACK client name
// including the final NULL character.  This value is a constant.
func ClientNameSize() int {
	return int(C.jack_client_name_size())
}

// GetClientPID returns the PID of the named client. If not available, 0 will be returned.
func GetClientPID(name string) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return int(C.jack_get_client_pid(cname))
}

//-----------------------------------------------------------------------------

// PortID is the Go type for the JACK port identifier.
type PortID uint32

// Port is the Go type for the JACK port structure.
type Port struct {
	ptr *C.struct__jack_port
}

//-----------------------------------------------------------------------------

// Client is the Go type for the JACK client structure.
type Client struct {
	ptr                        *C.struct__jack_client
	processCallback            funcProcessCallback
	bufferSizeCallback         funcBufferSizeCallback
	sampleRateCallback         funcSampleRateCallback
	portRegistrationCallback   funcPortRegistrationCallback
	portRenameCallback         funcPortRenameCallback
	portConnectCallback        funcPortConnectCallback
	clientRegistrationCallback funcClientRegistrationCallback
	freewheelCallback          funcFreewheelCallback
	graphOrderCallback         funcGraphOrderCallback
	xrunCallback               funcXrunCallback
	shutdownCallback           funcShutdownCallback
	infoShutdownCallback       funcInfoShutdownCallback
}

var (
	clientMap     map[*C.struct__jack_client]*Client
	clientMapLock sync.Mutex
	errorFunction funcErrorFunction
	infoFunction  funcInfoFunction
)

// ClientOpen opens an external client session with a JACK server.
func ClientOpen(name string, options int) (*Client, Status) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var status C.int
	cc := C.jack_client_open_go(cname, C.int(options), &status)
	var c *Client
	if cc != nil {
		if clientMap == nil {
			clientMap = make(map[*C.struct__jack_client]*Client)
		}
		c = new(Client)
		c.ptr = cc
		clientMapLock.Lock()
		clientMap[cc] = c
		clientMapLock.Unlock()
	}
	return c, Status(status)
}

// Close disconnects an external client from a JACK server.
func (c *Client) Close() int {
	if c == nil || c.ptr == nil {
		return 0
	}
	rc := int(C.jack_client_close(c.ptr))
	if rc == 0 {
		clientMapLock.Lock()
		delete(clientMap, c.ptr)
		clientMapLock.Unlock()
		c.ptr = nil
	}
	return rc
}

// GetName returns the actual client name.
func (c *Client) GetName() string {
	return C.GoString(C.jack_get_client_name(c.ptr))
}

// GetUUIDByName returns the session ID for a client name.
func (c *Client) GetUUIDByName(name string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return C.GoString(C.jack_get_uuid_for_client_name(c.ptr, cname))
}

// GetNameByUUID returns the client name for a session_id.
func (c *Client) GetNameByUUID(uuid string) string {
	cuuid := C.CString(uuid)
	defer C.free(unsafe.Pointer(cuuid))
	return C.GoString(C.jack_get_client_name_by_uuid(c.ptr, cuuid))
}

// Activate tells the Jack server that the program is ready to start processing audio.
func (c *Client) Activate() int {
	return int(C.jack_activate(c.ptr))
}

// Deactivate tells the Jack server to remove this client from the process graph.
func (c *Client) Deactivate() int {
	return int(C.jack_deactivate(c.ptr))
}

// OnShutdown registers a function to be called if the JACK server shuts down the client thread.
func (c *Client) OnShutdown(cb funcShutdownCallback) {
	c.shutdownCallback = cb
	C.jack_on_shutdown_go(c.ptr)
}

// OnInfoShutdown registers a function to be called if the JACK server shuts down the client thread.
func (c *Client) OnInfoShutdown(cb funcInfoShutdownCallback) {
	c.infoShutdownCallback = cb
	C.jack_on_info_shutdown_go(c.ptr)
}

// SetProcessCallback tells JACK to call the process callback whenever there is work be done.
func (c *Client) SetProcessCallback(cb funcProcessCallback) int {
	c.processCallback = cb
	return int(C.jack_set_process_callback_go(c.ptr))
}

// SetFreewheelCallback tells JACK to call the freewheel callback whenever we enter or leave "freewheel" mode.
func (c *Client) SetFreewheelCallback(cb funcFreewheelCallback) int {
	c.freewheelCallback = cb
	return int(C.jack_set_freewheel_callback_go(c.ptr))
}

// SetBufferSizeCallback tells JACK to call the bufsize callback whenever the size of the the buffer passed to the process callback is about to change.
func (c *Client) SetBufferSizeCallback(cb funcBufferSizeCallback) int {
	c.bufferSizeCallback = cb
	return int(C.jack_set_buffer_size_callback_go(c.ptr))
}

// SetSampleRateCallback tells JACK to call srate callback whenever the system sample rate changes.
func (c *Client) SetSampleRateCallback(cb funcSampleRateCallback) int {
	c.sampleRateCallback = cb
	return int(C.jack_set_sample_rate_callback_go(c.ptr))
}

// SetClientRegistrationCallback tells JACK to call client registration callback whenever a client is registered or unregistered.
func (c *Client) SetClientRegistrationCallback(cb funcClientRegistrationCallback) int {
	c.clientRegistrationCallback = cb
	return int(C.jack_set_client_registration_callback_go(c.ptr))
}

// SetPortRegistrationCallback tells JACK to call registration callback whenever a port is registered or unregistered.
func (c *Client) SetPortRegistrationCallback(cb funcPortRegistrationCallback) int {
	c.portRegistrationCallback = cb
	return int(C.jack_set_port_registration_callback_go(c.ptr))
}

// SetPortConnectCallback tells JACK to call connect callback whenever a port is connected or disconnected.
func (c *Client) SetPortConnectCallback(cb funcPortConnectCallback) int {
	c.portConnectCallback = cb
	return int(C.jack_set_port_connect_callback_go(c.ptr))
}

// SetPortRenameCallback tells JACK to call rename callback whenever a port is renamed.
func (c *Client) SetPortRenameCallback(cb funcPortRenameCallback) int {
	c.portRenameCallback = cb
	return int(C.jack_set_port_rename_callback_go(c.ptr))
}

// SetGraphOrderCallback tells JACK to call graph callback whenever the processing graph is reordered.
func (c *Client) SetGraphOrderCallback(cb funcGraphOrderCallback) int {
	c.graphOrderCallback = cb
	return int(C.jack_set_graph_order_callback_go(c.ptr))
}

// SetXrunCallback tells JACK to call xrun callback whenever there is an xrun.
func (c *Client) SetXrunCallback(cb funcXrunCallback) int {
	c.xrunCallback = cb
	return int(C.jack_set_xrun_callback_go(c.ptr))
}

// SetBufferSize changes the buffer size passed to the process callback.
func (c *Client) SetBufferSize(nframes uint32) int {
	return int(C.jack_set_buffer_size(c.ptr, C.jack_nframes_t(nframes)))
}

// GetSampleRate returns the sample rate of the jack system, as set by the user when jackd was started.
func (c *Client) GetSampleRate() uint32 {
	return uint32(C.jack_get_sample_rate(c.ptr))
}

// GetBufferSize returns the current maximum size that will ever be passed to the process callback.
func (c *Client) GetBufferSize() uint32 {
	return uint32(C.jack_get_buffer_size(c.ptr))
}

// PortRegister creates a new port for the client.
func (c *Client) PortRegister(portName, portType string, flags, bufferSize uint64) *Port {
	cportName := C.CString(portName)
	defer C.free(unsafe.Pointer(cportName))
	cportType := C.CString(portType)
	defer C.free(unsafe.Pointer(cportType))
	cport := C.jack_port_register(c.ptr, cportName, cportType, C.ulong(flags), C.ulong(bufferSize))
	if cport != nil {
		return &Port{cport}
	}
	return nil
}

// PortUnregister removes the port from the client, disconnecting any existing connections.
func (c *Client) PortUnregister(port *Port) int {
	return int(C.jack_port_unregister(c.ptr, port.ptr))
}

//-----------------------------------------------------------------------------

// int jack_set_latency_callback (jack_client_t *client,
// int jack_set_freewheel(jack_client_t* client, int onoff) JACK_OPTIONAL_WEAK_EXPORT;
// jack_native_thread_t jack_client_thread_id (jack_client_t *client) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_is_realtime (jack_client_t *client) JACK_OPTIONAL_WEAK_EXPORT;
// jack_nframes_t jack_thread_wait (jack_client_t *client, int status) JACK_OPTIONAL_WEAK_EXPORT;
// jack_nframes_t jack_cycle_wait (jack_client_t* client) JACK_OPTIONAL_WEAK_EXPORT;
// void jack_cycle_signal (jack_client_t* client, int status) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_set_process_thread(jack_client_t* client, JackThreadCallback thread_callback, void *arg) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_set_thread_init_callback (jack_client_t *client,

// float jack_cpu_load (jack_client_t *client) JACK_OPTIONAL_WEAK_EXPORT;
// void * jack_port_get_buffer (jack_port_t *port, jack_nframes_t) JACK_OPTIONAL_WEAK_EXPORT;
// jack_uuid_t jack_port_uuid (const jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// const char * jack_port_name (const jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// const char * jack_port_short_name (const jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_flags (const jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// const char * jack_port_type (const jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// jack_port_type_id_t jack_port_type_id (const jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_is_mine (const jack_client_t *client, const jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_connected (const jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_connected_to (const jack_port_t *port,
// const char ** jack_port_get_connections (const jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// const char ** jack_port_get_all_connections (const jack_client_t *client,
// int jack_port_rename (jack_client_t* client, jack_port_t *port, const char *port_name) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_set_alias (jack_port_t *port, const char *alias) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_unset_alias (jack_port_t *port, const char *alias) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_get_aliases (const jack_port_t *port, char* const aliases[2]) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_request_monitor (jack_port_t *port, int onoff) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_request_monitor_by_name (jack_client_t *client,
// int jack_port_ensure_monitor (jack_port_t *port, int onoff) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_monitoring_input (jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_connect (jack_client_t *client,
// int jack_disconnect (jack_client_t *client,
// int jack_port_disconnect (jack_client_t *client, jack_port_t *port) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_name_size(void) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_port_type_size(void) JACK_OPTIONAL_WEAK_EXPORT;
// size_t jack_port_type_get_buffer_size (jack_client_t *client, const char *port_type) JACK_WEAK_EXPORT;
// void jack_port_get_latency_range (jack_port_t *port, jack_latency_callback_mode_t mode, jack_latency_range_t *range) JACK_WEAK_EXPORT;
// void jack_port_set_latency_range (jack_port_t *port, jack_latency_callback_mode_t mode, jack_latency_range_t *range) JACK_WEAK_EXPORT;
// int jack_recompute_total_latencies (jack_client_t *client) JACK_OPTIONAL_WEAK_EXPORT;
// jack_nframes_t jack_port_get_total_latency (jack_client_t *client,
// const char ** jack_get_ports (jack_client_t *client,
// jack_port_t * jack_port_by_name (jack_client_t *client, const char *port_name) JACK_OPTIONAL_WEAK_EXPORT;
// jack_port_t * jack_port_by_id (jack_client_t *client,
// jack_nframes_t jack_frames_since_cycle_start (const jack_client_t *) JACK_OPTIONAL_WEAK_EXPORT;
// jack_nframes_t jack_frame_time (const jack_client_t *) JACK_OPTIONAL_WEAK_EXPORT;
// jack_nframes_t jack_last_frame_time (const jack_client_t *client) JACK_OPTIONAL_WEAK_EXPORT;
// int jack_get_cycle_times(const jack_client_t *client,
// jack_time_t jack_frames_to_time(const jack_client_t *client, jack_nframes_t) JACK_OPTIONAL_WEAK_EXPORT;
// jack_nframes_t jack_time_to_frames(const jack_client_t *client, jack_time_t) JACK_OPTIONAL_WEAK_EXPORT;
// jack_time_t jack_get_time(void) JACK_OPTIONAL_WEAK_EXPORT;
// extern void (*jack_error_callback)(const char *msg) JACK_OPTIONAL_WEAK_EXPORT;
// void jack_set_error_function (void (*func)(const char *)) JACK_OPTIONAL_WEAK_EXPORT;
// extern void (*jack_info_callback)(const char *msg) JACK_OPTIONAL_WEAK_EXPORT;
// void jack_set_info_function (void (*func)(const char *)) JACK_OPTIONAL_WEAK_EXPORT;
// void jack_free(void* ptr) JACK_OPTIONAL_WEAK_EXPORT;

//-----------------------------------------------------------------------------
