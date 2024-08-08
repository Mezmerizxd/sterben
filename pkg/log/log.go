package log

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"github.com/rs/zerolog"
)

// Default directory and file mode for log files.
var logsDirectory = "logs/"
var logsFileMode fs.FileMode = 0755

// Config struct holds the configuration options for the logger.
type Config struct {
	Feature      string `json:"feature"`       // The name of the feature being logged.
	ConsoleOutput bool  `json:"consoleOutput"` // Flag to enable/disable console output.
	FileOutput    bool  `json:"fileOutput"`    // Flag to enable/disable file output.
}

// Log interface defines the methods available for logging.
type Log interface {
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Debug() *zerolog.Event
	Close() error
}

// log struct implements the Log interface and holds logger configuration and state.
type log struct {
	config        Config
	writer        zerolog.Logger
	file          *os.File
	consoleWriter *ConsoleWriter
}

// New creates a new logger instance based on the provided configuration.
func New(cfg Config) Log {
	var writers []io.Writer
	var file *os.File

	// Check if file output is enabled.
	if cfg.FileOutput {
		// Ensure the logs directory exists.
		if _, err := os.Stat(logsDirectory); os.IsNotExist(err) {
			if err := os.Mkdir(logsDirectory, logsFileMode); err != nil {
				print("Error creating logs directory")
				panic(err)
			}
		}

		// Open or create the log file.
		var err error
		file, err = os.OpenFile(logsDirectory+cfg.Feature+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, logsFileMode)
		if err != nil {
			print("Error opening log file")
			panic(err)
		}
		writers = append(writers, file)
	}

	// If console output is enabled, add os.Stdout to the writers.
	var consoleWriter *ConsoleWriter
	if cfg.ConsoleOutput {
		consoleWriter = NewConsoleWriter()
		writers = append(writers, consoleWriter)
	}

	// Create a multi-writer to output to both file and console if both are enabled.
	multiWriter := io.MultiWriter(writers...)

	// Create a zerolog logger using the multi-writer.
	writer := zerolog.New(multiWriter).With().Timestamp().Logger()

	return &log{
		config:        cfg,
		writer:        writer,
		file:          file,
		consoleWriter: consoleWriter,
	}
}

// Info logs an informational message.
func (l *log) Info() *zerolog.Event {
	return l.writer.Info().Str("feature", l.config.Feature)
}

// Warn logs a warning message.
func (l *log) Warn() *zerolog.Event {
	return l.writer.Warn().Str("feature", l.config.Feature)
}

// Error logs an error message.
func (l *log) Error() *zerolog.Event {
	return l.writer.Error().Str("feature", l.config.Feature)
}

// Debug logs a debug message.
func (l *log) Debug() *zerolog.Event {
	return l.writer.Debug().Str("feature", l.config.Feature)
}

// Close finalizes the logging by closing any open resources, such as file handles.
func (l *log) Close() error {
	l.Info().Msgf("Shutting down %s", l.config.Feature)

	// Close the file only if file output is enabled and the file is not nil.
	if l.config.FileOutput && l.file != nil {
		l.file.Close()
	}

	return nil
}

// WriterOutputFormat defines the structure of a log output format.
type WriterOutputFormat struct {
	Level   string `json:"level"`   // Log level (info, warn, error, debug).
	Feature string `json:"feature"` // The feature associated with the log message.
	Time    string `json:"time"`    // Timestamp of the log event.
	Message string `json:"message"` // The actual log message.
}

// GetWriterOutputFormat parses JSON log data into the WriterOutputFormat struct.
func GetWriterOutputFormat(data string) WriterOutputFormat {
	var output WriterOutputFormat
	err := json.Unmarshal([]byte(data), &output)
	if err != nil {
		print("Error unmarshalling log data")
		panic(err)
	}
	return output
}
