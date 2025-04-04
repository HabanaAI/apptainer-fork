// Copyright (c) Contributors to the Apptainer project, established as
//   Apptainer a Series of LF Projects LLC.
//   For website terms of use, trademark policy, privacy policy and other
//   project policies see https://lfprojects.org/policies
// Copyright (c) 2018-2024, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package sources

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apptainer/apptainer/pkg/build/types"
)

// ScratchConveyor only needs to hold the conveyor to have the needed data to pack
type ScratchConveyor struct {
	b *types.Bundle
}

// ScratchConveyorPacker only needs to hold the conveyor to have the needed data to pack
type ScratchConveyorPacker struct {
	ScratchConveyor
}

// Get just stores the source
func (c *ScratchConveyor) Get(_ context.Context, b *types.Bundle) (err error) {
	c.b = b

	return nil
}

// Pack puts relevant objects in a Bundle!
func (cp *ScratchConveyorPacker) Pack(context.Context) (b *types.Bundle, err error) {
	err = cp.insertBaseEnv()
	if err != nil {
		return nil, fmt.Errorf("while inserting base environment: %v", err)
	}

	err = cp.insertRunScript()
	if err != nil {
		return nil, fmt.Errorf("while inserting runscript: %v", err)
	}

	return cp.b, nil
}

func (c *ScratchConveyor) insertBaseEnv() (err error) {
	if err = makeBaseEnv(c.b.RootfsPath, true); err != nil {
		return
	}
	return nil
}

func (cp *ScratchConveyorPacker) insertRunScript() (err error) {
	err = os.WriteFile(filepath.Join(cp.b.RootfsPath, "/.singularity.d/runscript"), []byte("#!/bin/sh\n"), 0o755)
	if err != nil {
		return
	}

	return nil
}

// CleanUp removes any tmpfs owned by the conveyorPacker on the filesystem
func (c *ScratchConveyor) CleanUp() {
	c.b.Remove()
}
