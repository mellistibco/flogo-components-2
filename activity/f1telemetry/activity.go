package f1telemetry

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/lunixbochs/struc"
	//"gopkg.in/restruct.v1"
)

// activityLog is the default logger for the Log Activity
var log = logger.GetLogger("activity-f1telemetry")

const (
	ivInput = "buffer"

	ovOutput  = "cardata"
	ovOutput2 = "cararray"
)

// F1Data - Struct for the unpacking of the UDP data format
type F1Data struct {
	Time                 float32 `struc:"float32,little"` //F
	LapTime              float32 `struc:"float32,little"` //F
	LapDistance          float32 `struc:"float32,little"` //F
	TotalDistance        float32 `struc:"float32,little"` //F
	X                    float32 `struc:"float32,little"` // World space position F
	Y                    float32 `struc:"float32,little"` // World space position F
	Z                    float32 `struc:"float32,little"` // World space position F
	Speed                float32 `struc:"float32,little"` // Speed of car in MPH F
	Xv                   float32 `struc:"float32,little"` // Velocity in world space F
	Yv                   float32 `struc:"float32,little"` // Velocity in world space F
	Zv                   float32 `struc:"float32,little"` // Velocity in world space F
	Xr                   float32 `struc:"float32,little"` // World space right direction F
	Yr                   float32 `struc:"float32,little"` // World space right direction F
	Zr                   float32 `struc:"float32,little"` // World space right direction F
	Xd                   float32 `struc:"float32,little"` // World space forward direction F
	Yd                   float32 `struc:"float32,little"` // World space forward direction F
	Zd                   float32 `struc:"float32,little"` // World space forward direction F
	SuspPosRL            float32 `struc:"float32,little"` // Suspension position RL, RR, FL, FR  F array
	SuspPosRR            float32 `struc:"float32,little"`
	SuspPosFL            float32 `struc:"float32,little"`
	SuspPosFR            float32 `struc:"float32,little"`
	SuspVelRL            float32 `struc:"float32,little"` // Suspension velocity RL, RR, FL, FR  F array
	SuspVelRR            float32 `struc:"float32,little"`
	SuspVelFL            float32 `struc:"float32,little"`
	SuspVelFR            float32 `struc:"float32,little"`
	WheelspeedRL         float32 `struc:"float32,little"` // Wheel Speed RL, RR, FL, FR  F array
	WheelspeedRR         float32 `struc:"float32,little"`
	WheelspeedFL         float32 `struc:"float32,little"`
	WheelspeedFR         float32 `struc:"float32,little"`
	Throttle             float32 `struc:"float32,little"` // F
	Steer                float32 `struc:"float32,little"` // F
	Brake                float32 `struc:"float32,little"` // F
	Clutch               float32 `struc:"float32,little"` // F
	Gear                 float32 `struc:"float32,little"` // F
	Gforcelat            float32 `struc:"float32,little"` // F
	Gforcelon            float32 `struc:"float32,little"` // F
	Lap                  float32 `struc:"float32,little"` // F
	EngineRate           float32 `struc:"float32,little"` // F
	Slipronativesupport  float32 `struc:"float32,little"` // F	// SLI Pro support
	Carposition          float32 `struc:"float32,little"` // F	// car race position
	Kerslevel            float32 `struc:"float32,little"` // F	// kers energy left
	Kersmaxlevel         float32 `struc:"float32,little"` // F	// kers maximum energy
	Drs                  float32 `struc:"float32,little"` // F	// 0 = off, 1 = on
	Tractioncontrol      float32 `struc:"float32,little"` // F	// 0 (off) - 2 (high)
	Antilockbrakes       float32 `struc:"float32,little"` // F	// 0 (off) - 1 (on)
	Fuelintank           float32 `struc:"float32,little"` // F	// current fuel mass
	Fuelcapacity         float32 `struc:"float32,little"` // F	// fuel capacity
	Inpits               float32 `struc:"float32,little"` // F	// 0 = none, 1 = pitting, 2 = in pit area
	Sector               float32 `struc:"float32,little"` // F	// 0 = sector1, 1 = sector2, 2 = sector3
	Sector1time          float32 `struc:"float32,little"` // F	// time of sector1 (or 0)
	Sector2time          float32 `struc:"float32,little"` // F	// time of sector2 (or 0)
	BrakestempRL         float32 `struc:"float32,little"` // brakes temperature (centigrade)
	BrakestempRR         float32 `struc:"float32,little"`
	BrakestempFL         float32 `struc:"float32,little"`
	BrakestempFR         float32 `struc:"float32,little"`
	TyrespressureRL      float32 `struc:"float32,little"` // tyres pressure PSI
	TyrespressureRR      float32 `struc:"float32,little"`
	TyrespressureFL      float32 `struc:"float32,little"`
	TyrespressureFR      float32 `struc:"float32,little"`
	Teaminfo             float32 `struc:"float32,little"` // F	// team ID
	Totallaps            float32 `struc:"float32,little"` // F	// total number of laps in this race
	Tracksize            float32 `struc:"float32,little"` // F	// track size meters
	Lastlaptime          float32 `struc:"float32,little"` // F	// last lap time
	Maxrpm               float32 `struc:"float32,little"` // cars max RPM, at which point the rev limiter will kick in
	Idlerpm              float32 `struc:"float32,little"` // cars idle RPM
	Maxgears             float32 `struc:"float32,little"` // maximum number of gears
	SessionType          float32 `struc:"float32,little"` // 0 = unknown, 1 = practice, 2 = qualifying, 3 = race
	DrsAllowed           float32 `struc:"float32,little"` // 0 = not allowed, 1 = allowed, -1 = invalid / unknown
	Tracknumber          float32 `struc:"float32,little"` // -1 for unknown, 0-21 for tracks
	VehicleFIAFlags      float32 `struc:"float32,little"` // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
	Era                  float32 `struc:"float32,little"` // era, 2017 (modern) or 1980 (classic)
	Enginetemperature    float32 `struc:"float32,little"` // engine temperature (centigrade)
	Gforcevert           float32 `struc:"float32,little"` // vertical g-force component
	Angvelx              float32 `struc:"float32,little"` // angular velocity x-component
	Angvely              float32 `struc:"float32,little"` // angular velocity y-component
	Angvelz              float32 `struc:"float32,little"` // angular velocity z-component
	TyrestemperatureRL   byte    `struc:"byte,little"`    // tyres temperature (centigrade)
	TyrestemperatureRR   byte    `struc:"byte,little"`
	TyrestemperatureFL   byte    `struc:"byte,little"`
	TyrestemperatureFR   byte    `struc:"byte,little"`
	TyreswearRL          byte    `struc:"byte,little"` // tyre wear percentage
	TyreswearRR          byte    `struc:"byte,little"`
	TyreswearFL          byte    `struc:"byte,little"`
	TyreswearFR          byte    `struc:"byte,little"`
	Tyrecompound         byte    `struc:"byte,little"` // compound of tyre – 0 = ultra soft, 1 = super soft, 2 = soft, 3 = medium, 4 = hard, 5 = inter, 6 = wet
	Frontbrakebias       byte    `struc:"byte,little"` // front brake bias (percentage)
	Fuelmix              byte    `struc:"byte,little"` // fuel mix - 0 = lean, 1 = standard, 2 = rich, 3 = max
	CurrentLapInvalid    byte    `struc:"byte,little"` // current lap invalid - 0 = valid, 1 = invalid
	TyresdamageRL        byte    `struc:"byte,little"` // tyre damage (percentage)
	TyresdamageRR        byte    `struc:"byte,little"`
	TyresdamageFL        byte    `struc:"byte,little"`
	TyresdamageFR        byte    `struc:"byte,little"`
	Frontleftwingdamage  byte    `struc:"byte,little"`    // front left wing damage (percentage)
	Frontrightwingdamage byte    `struc:"byte,little"`    // front right wing damage (percentage)
	Rearwingdamage       byte    `struc:"byte,little"`    // rear wing damage (percentage)
	Enginedamage         byte    `struc:"byte,little"`    // engine damage (percentage)
	Gearboxdamage        byte    `struc:"byte,little"`    // gear box damage (percentage)
	Exhaustdamage        byte    `struc:"byte,little"`    // exhaust damage (percentage)
	Pitlimiterstatus     byte    `struc:"byte,little"`    // pit limiter status – 0 = off, 1 = on
	Pitspeedlimit        byte    `struc:"byte,little"`    // pit speed limit in mph
	Sessiontimeleft      float32 `struc:"float32,little"` // NEW: time left in session in seconds
	Revlightspercent     byte    `struc:"byte,little"`    // NEW: rev lights indicator (percentage)
	Isspectating         byte    `struc:"byte,little"`    // NEW: whether the player is spectating
	Spectatorcarindex    byte    `struc:"byte,little"`    // NEW: index of the car being spectated
	NumCars              byte    `struc:"byte,little"`    // number of cars in data
	PlayerCarIndex       byte    `struc:"byte,little"`
	Filler1              []byte  `struc:"[900]byte"`      // cars data array
	Yaw                  float32 `struc:"float32,little"` // NEW (v1.8)
	Pitch                float32 `struc:"float32,little"` // NEW (v1.8)
	Roll                 float32 `struc:"float32,little"` // NEW (v1.8)
	XLocalVelocity       float32 `struc:"float32,little"` // NEW (v1.8) Velocity in local space
	YLocalVelocity       float32 `struc:"float32,little"` // NEW (v1.8) Velocity in local space
	ZLocalVelocity       float32 `struc:"float32,little"` // NEW (v1.8) Velocity in local space
	SuspAccelerationRL   float32 `struc:"float32,little"` // NEW (v1.8) RL, RR, FL, FR
	SuspAccelerationRR   float32 `struc:"float32,little"`
	SuspAccelerationFL   float32 `struc:"float32,little"`
	SuspAccelerationFR   float32 `struc:"float32,little"`
	AngAccX              float32 `struc:"float32,little"` // NEW (v1.8) angular acceleration x-component
	AngAccY              float32 `struc:"float32,little"` // NEW (v1.8) angular acceleration y-component
	AngAccZ              float32 `struc:"float32,little"` // NEW (v1.8) angular acceleration z-component
}

// F1CarArray - Struct for the unpacking of the UDP data format (Car data array)
type F1CarArray struct {
	X                 float32 `struc:"float32,little"` // world co-ordinates of vehicle
	Y                 float32 `struc:"float32,little"` // world co-ordinates of vehicle
	Z                 float32 `struc:"float32,little"` // world co-ordinates of vehicle
	LastLapTime       float32 `struc:"float32,little"`
	CurrentLapTime    float32 `struc:"float32,little"`
	BestLapTime       float32 `struc:"float32,little"`
	Sector1Time       float32 `struc:"float32,little"`
	Sector2Time       float32 `struc:"float32,little"`
	LapDistance       float32 `struc:"float32,little"`
	DriverID          byte    `struc:"byte"`
	TeamID            byte    `struc:"byte"`
	CarPosition       byte    `struc:"byte"` // UPDATED: track positions of vehicle
	CurrentLapNum     byte    `struc:"byte"`
	TyreCompound      byte    `struc:"byte"` // compound of tyre – 0 = ultra soft, 1 = super soft, 2 = soft, 3 = medium, 4 = hard, 5 = inter, 6 = wet
	InPits            byte    `struc:"byte"` // 0 = none, 1 = pitting, 2 = in pit area
	Sector            byte    `struc:"byte"` // 0 = sector1, 1 = sector2, 2 = sector3
	CurrentLapInvalid byte    `struc:"byte"` // current lap invalid - 0 = valid, 1 = invalid
	Penalties         byte    `struc:"byte"` // NEW: accumulated time penalties in seconds to be added
}

func init() {
	log.SetLogLevel(logger.InfoLevel)
	//log.SetLogLevel(logger.DebugLevel)
}

// f1telemetry is an Activity that takes in data from a byte stream and interprets it as data from F1 2017
//
// inputs : {buffer} (byte array) RAW UDP data
// outputs: {data} (string) CSV data
type f1telemetry struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &f1telemetry{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *f1telemetry) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *f1telemetry) Eval(context activity.Context) (done bool, err error) {

	// Get the runtime values
	log.Debug("Starting")

	input, _ := context.GetInput(ivInput).([]byte)
	buf := bytes.NewBuffer(input)

	log.Debugf("input : \n %x \n", input)

	// Create structs to hold unpacked data
	unpackedData := &F1Data{}
	unpackedData2 := &F1CarArray{}

	log.Debug("Unpack")

	// Unpack the main data
	err = struc.Unpack(buf, unpackedData)
	if err != nil {
		log.Error("Unpack Fail: F1Data ", err.Error())
		return false, err
	}

	log.Debug("print")
	log.Debugf("struct F1Data : \n %+v \n", unpackedData)

	// Write the CSV rows to the output
	fields := unpackedData.valueStrings()
	fieldsstring := strings.Join(fields, ",")

	log.Debugf("CSV data : %v \n", fieldsstring)
	context.SetOutput(ovOutput, fieldsstring)

	// Unpack the 20 item car data array
	// First load the 900 bytes into a bytes buffer and initialise the output "CSV" header
	// Note - Output array is:
	// Timestamp + array of car CSV data seprated by a "|"

	buf2 := bytes.NewBuffer(unpackedData.Filler1)
	arraystring := ""

	for i := 0; i <= 19; i++ {
		err = struc.Unpack(buf2, unpackedData2)
		if err != nil {
			log.Error("Unpack Fail: F1CarArray ", err.Error())
			return false, err
		}
		log.Debugf("Car Array unpacked: %v\n%+v\n", i, unpackedData2)
		arrayfields := unpackedData2.valueStrings()
		if arraystring == "" {
			arraystring = fmt.Sprintf("%v", i) + "," + fmt.Sprintf("%g", unpackedData.Time) + "," + strings.Join(arrayfields, ",")
		} else {
			arraystring = arraystring + "|" + fmt.Sprintf("%v", i) + "," + fmt.Sprintf("%g", unpackedData.Time) + "," + strings.Join(arrayfields, ",")
		}

	}

	// Write the CSV array rows to the output
	log.Debugf("CSV Car Array data : %v \n", arraystring)
	context.SetOutput(ovOutput2, arraystring)

	return true, nil
}
func (f F1Data) valueStrings() []string {
	v := reflect.ValueOf(f)
	ss := make([]string, v.NumField())
	for i := range ss {
		typeField := v.Type().Field(i)
		if strings.HasPrefix(typeField.Name, "Filler") {
			ss[i] = fmt.Sprintf("%v", "-1")
		} else {
			switch v.Field(i).Kind() {
			case reflect.Float32, reflect.Float64:
				x := v.Field(i).Float()
				ss[i] = strconv.FormatFloat(x, 'f', -1, 32)
			default:
				ss[i] = fmt.Sprintf("%v", v.Field(i))
			}

		}
	}
	return ss
}
func (f F1CarArray) valueStrings() []string {
	v := reflect.ValueOf(f)
	ss := make([]string, v.NumField())
	for i := range ss {
		//typeField := v.Type().Field(i)
		switch v.Field(i).Kind() {
		case reflect.Float32, reflect.Float64:
			x := v.Field(i).Float()
			ss[i] = strconv.FormatFloat(x, 'f', -1, 32)
		default:
			ss[i] = fmt.Sprintf("%v", v.Field(i))
		}

	}
	return ss
}
