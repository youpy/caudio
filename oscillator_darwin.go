package caudio

/*
#cgo LDFLAGS: -framework AudioUnit
#include <AudioToolbox/AudioToolbox.h>
#include  "callback.h"

typedef struct synthDef {
  UInt32 frameCount;
  int callbackIndex;
  int stepCount;
  float lastValue;
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

  // TODO: make fps configurable
  int framesPerStep = 5;

  float timeInSeconds;
  float value;

  float *outL = ioData->mBuffers[0].mData;
  float *outR = ioData->mBuffers[1].mData;

  int i;
  for (i=0; i< inNumberFrames; i++){
    if (def->frameCount % framesPerStep == 0) {
      timeInSeconds = def->frameCount / samplingRate;
      value = baudio_callback(def->callbackIndex, timeInSeconds, def->stepCount);

      def->lastValue = value;
      def->stepCount = def->stepCount + 1;
    }

    *outL++ = def->lastValue;
    *outR++ = def->lastValue;

    def->frameCount = def->frameCount + 1;
  }

  return noErr;
}

synthDef *createSynthDef(int callbackIndex)
{
  synthDef *def = malloc(sizeof(synthDef));

  def->frameCount = 0;
  def->stepCount = 0;
  def->callbackIndex = callbackIndex;
  def->lastValue = 0.0;

  return def;
}
*/
import "C"
import (
	"unsafe"
)

type Oscillator struct {
	synthDef *C.struct_synthDef
}

func New(fn Callback) *Oscillator {
	callbackIndex := register(fn)
	synthDef := C.createSynthDef(C.int(callbackIndex))

	return &Oscillator{synthDef}
}

func (o *Oscillator) Start() error {
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
