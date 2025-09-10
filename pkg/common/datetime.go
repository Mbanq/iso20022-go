package common

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

// ISODate is a custom type wrapping time.Time
type ISODate time.Time

func (i ISODate) MarshalText() (string, error) {
	return time.Time(i).Format("2006-01-02"), nil
}

func (i ISODate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	content, err := i.MarshalText()
	if err != nil {
		return err
	}
	return e.EncodeElement([]byte(content), start)
}

func (i *ISODate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := d.DecodeElement(&content, &start); err != nil {
		return err
	}

	// Parse the date in YYYY-MM-DD format
	tt, err := time.Parse("2006-01-02", content)
	if err != nil {
		return fmt.Errorf("invalid ISODate format: %v", err)
	}

	*i = ISODate(tt)
	return nil
}

func (i *ISODate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ISODate should be a string, got %s", data)
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("failed to parse ISODate: %w", err)
	}
	*i = ISODate(t)
	return nil
}

func (i ISODate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(i).Format("2006-01-02"))
}

type ISODateTime time.Time

func (t *ISODateTime) UnmarshalJSON(data []byte) error {
	// ISODateTime should be a string inside quotes
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ISODateTime should be a string, got %s", data)
	}

	// Parse the time string using RFC3339 format
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("failed to parse ISODateTime string '%s': %w", s, err)
	}

	*t = ISODateTime(parsedTime)
	return nil
}

func (t ISODateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.RFC3339))
}

func (t ISODateTime) Validate() error {
	return nil
}

func (t *ISODateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return err
	}
	*t = ISODateTime(parse)
	return nil
}

func (t ISODateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(time.Time(t).Format(time.RFC3339), start)
}

// ISOTime is a custom type wrapping time.Time for time-only values
type ISOTime time.Time

func (t *ISOTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ISOTime should be a string, got %s", data)
	}
	// ISO 8601 time format is 15:04:05Z07:00, we'll use a common representation
	parsedTime, err := time.Parse("15:04:05", s)
	if err != nil {
		return fmt.Errorf("failed to parse ISOTime: %w", err)
	}
	*t = ISOTime(parsedTime)
	return nil
}

func (t ISOTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("15:04:05"))
}
