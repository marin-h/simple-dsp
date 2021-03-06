package test

import (
	"testing"
	"time"

	"github.com/marin-h/simple-dsp/app"
	cmd "github.com/marin-h/simple-dsp/cmd"
)

func TestAPIFrequencyCapped(t *testing.T) {

	var i int64
	for i = 0; i < cmd.Dsp.MaxImpressionsPerMinute-1; i++ {
		now := time.Now()
		if cmd.Dsp.FrequencyCapped("a", now) {
			t.Errorf("User should not be capped at the moment!")
		} else {
			t.Logf("User is not capped at %s", now)
		}
		TestHandleBid(t)
		TestWinNotice(t)
	}

	if !cmd.Dsp.FrequencyCapped("a", time.Now()) {
		t.Errorf("User should by capped by now!")
	}
}

func TestTimestampsFrequencyCapping(t *testing.T) {

	var i int64
	now := time.Now()
	offset := -time.Second

	for i = 0; i < cmd.Dsp.MaxImpressionsPerMinute-1; i++ {

		cmd.Dsp.RegisterImpression(app.Bid{
			Id:        app.UUID(),
			UserId:    "kjkw340r",
			Timestamp: now.Add(offset).Unix(),
			Status:    "",
			Price:     0.1})

		t.Logf("Setting impression timestamp %s before", offset)
		t.Logf(">> adding impression %d", i+1)

		// Not yet capped
		if cmd.Dsp.FrequencyCapped("kjkw340r", now.Add(offset)) {
			t.Errorf("User is capped!")
		} else {
			t.Logf("User is not capped %s", offset)
		}
		offset = offset - time.Second*12
	}

	cmd.Dsp.RegisterImpression(app.Bid{
		Id:        app.UUID(),
		UserId:    "kjkw340r",
		Timestamp: now.Add(offset).Unix(),
		Status:    "",
		Price:     0.1})

	t.Logf("Setting impression timestamp %s before", offset)
	t.Logf(">> adding impression %d", i+1)

	// Just capped
	if !cmd.Dsp.FrequencyCapped("kjkw340r", now.Add(offset)) {
		t.Errorf("User is not capped %s", offset)
	} else {
		t.Logf("User is capped %s", offset)
	}

	// Capped some time after
	offset = time.Second * 11
	if !cmd.Dsp.FrequencyCapped("kjkw340r", now.Add(offset)) {
		t.Errorf("User is not capped %s", offset)
	} else {
		t.Logf("User is capped at %s", offset)
	}

	// Not capped anymore
	offset = time.Second * 12
	if cmd.Dsp.FrequencyCapped("kjkw340r", now.Add(offset)) {
		t.Errorf("User is capped! at %s", offset)
	} else {
		t.Logf("User is not capped %s", offset)
	}
}
