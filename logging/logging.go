// Copyright 2013, Cong Ding. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: Cong Ding <dinggnu@gmail.com>
//
package logging

import (
	"io"
	"os"
	"sync"
	"time"
)

// pre-defined formats
const (
	defaultFileName = "logging.log"
	configFileName  = "logging.conf"
)

// the logging struct
type Logger struct {
	// Be careful of the alignment issue of the variable seqid because it
	// uses the sync/atomic.AddUint64() operation. If the alignment is
	// wrong, it will cause a panic. To solve the alignment issue in an
	// easy way, we put seqid to the beginning of the structure.
	seqid     uint64
	name      string
	level     Level
	format    string
	out       io.Writer
	lock      sync.Mutex
	startTime time.Time
	sync      bool
}

// SimpleLogger creates a new logger with simple configuration.
func SimpleLogger(name string) *Logger {
	return createLogger(name, WARNING, BasicFormat, os.Stdout, true)
}

// BasieLogger creates a new logger with basic configuration.
func BasicLogger(name string) *Logger {
	return SimpleLogger(name)
}

// RichLogger creates a new logger with simple configuration.
func RichLogger(name string) *Logger {
	return FileLogger(name, NOTSET, RichFormat, defaultFileName, true)
}

// FileLogger creates a new logger with file output.
func FileLogger(name string, level Level, format string, file string, sync bool) *Logger {
	out, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	return createLogger(name, level, format, out, sync)
}

// createLogger create a new logger
func createLogger(name string, level Level, format string, out io.Writer, sync bool) *Logger {
	logger := new(Logger)
	logger.name = name
	logger.level = level
	logger.format = format
	logger.out = out
	logger.seqid = 0
	logger.sync = sync

	logger.init()
	return logger
}

// initialize the logger
func (logger *Logger) init() {
	logger.startTime = time.Now()
}

// get and set the configuration of the logger

func (logger *Logger) Name() string {
	return logger.name
}

func (logger *Logger) SetName(name string) {
	logger.name = name
}

func (logger *Logger) Level() Level {
	return logger.level
}

func (logger *Logger) SetLevel(level Level) {
	logger.level = Level(level)
}

func (logger *Logger) LevelName() string {
	name, _ := levelNames[logger.level]
	return name
}

func (logger *Logger) SetLevelName(name string) {
	level, ok := levelValues[name]
	if ok {
		logger.level = level
	}
}

func (logger *Logger) Format() string {
	return logger.format
}

func (logger *Logger) SetFormat(format string) {
	logger.format = format
}

func (logger *Logger) Writer() io.Writer {
	return logger.out
}

func (logger *Logger) SetWriter(out io.Writer) {
	logger.out = out
}
