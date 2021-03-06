package caudio

/*
#cgo LDFLAGS: -framework AudioUnit
#include <AudioToolbox/AudioToolbox.h>
#include  "callback.h"

typedef struct synthDef {
  UInt32 frameCount;
  int callbackIndex;
  int stepCount;
} synthDef;

const int sizeofAURenderCalalbackStruct = sizeof(AURenderCallbackStruct);

OSStatus RenderCallback(
  void                        *inRefCon,
  AudioUnitRenderActionFlags  *ioActionFlags,
  const AudioTimeStamp        *inTimeStamp,
  UInt32                      inBusNumber,
  UInt32                      inNumberFrames,
  AudioBufferList             *ioData)
{
  synthDef* def = inRefCon;
  float samplingRate = 44100;
  float timeInSeconds;
  float value;

  float *outL = ioData->mBuffers[0].mData;
  float *outR = ioData->mBuffers[1].mData;

  int i;
  for (i=0; i< inNumberFrames; i++){
    timeInSeconds = def->frameCount / samplingRate;
    value = baudio_callback(def->callbackIndex, timeInSeconds, def->stepCount);

    *outL++ = value;
    *outR++ = value;

    def->frameCount = def->frameCount + 1;
    def->stepCount = def->stepCount + 1;
  }

  return noErr;
}

synthDef *createSynthDef(int callbackIndex)
{
  synthDef *def = malloc(sizeof(synthDef));

  def->frameCount = 0;
  def->stepCount = 0;
  def->callbackIndex = callbackIndex;

  return def;
}
*/
import "C"
import (
	"unsafe"
)

// Audio is responsible for producing the sound using a given callback
type Audio struct {
	synthDef *C.struct_synthDef
}

func _new(fn Callback) *Audio {
	callbackIndex := register(fn)
	synthDef := C.createSynthDef(C.int(callbackIndex))

	return &Audio{synthDef}
}

func (o *Audio) _start() error {
	var defaultOutputUnit C.AudioUnit
	var cd C.AudioComponentDescription

	cd.componentType = C.kAudioUnitType_Output
	cd.componentSubType = C.kAudioUnitSubType_DefaultOutput
	cd.componentManufacturer = C.kAudioUnitManufacturer_Apple
	cd.componentFlags = 0
	cd.componentFlagsMask = 0

	comp := C.AudioComponentFindNext(nil, &cd)
	C.AudioComponentInstanceNew(comp, &defaultOutputUnit)

	var input C.AURenderCallbackStruct

	input.inputProc = C.AURenderCallback(C.RenderCallback)
	input.inputProcRefCon = unsafe.Pointer(o.synthDef)

	C.AudioUnitSetProperty(
		defaultOutputUnit,
		C.kAudioUnitProperty_SetRenderCallback,
		C.kAudioUnitScope_Input,
		0,
		unsafe.Pointer(&input),
		C.uint(unsafe.Sizeof(input)),
	)

	C.AudioUnitInitialize(defaultOutputUnit)
	C.AudioOutputUnitStart(defaultOutputUnit)

	return nil
}
