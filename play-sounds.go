package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var isInitialized bool

func initSpeaker(sampleRate beep.SampleRate) {
	if !isInitialized {
		speaker.Init(sampleRate, sampleRate.N(time.Second/10))
		isInitialized = true
	}
}

func playSound(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	initSpeaker(format.SampleRate) // 初回のみ初期化

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	<-done
	return nil
}
