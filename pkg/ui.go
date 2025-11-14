package pkg

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/fatih/color"
)

// FormatAndPrintResponse prints a structured APIResponse with nice formatting and colors.
func FormatAndPrintResponse(req APIRequest, resp APIResponse) {
    // Save request to history (best-effort)
    _ = SaveRequestHistory(req)

    // Error handling
    if resp.Error != "" {
        color.New(color.FgRed).Printf("Error: %s\n", resp.Error)
        return
    }

    // Status line: colorize based on status code
    statusColor := color.New(color.FgGreen)
    if resp.StatusCode >= 400 && resp.StatusCode < 500 {
        statusColor = color.New(color.FgYellow)
    } else if resp.StatusCode >= 500 {
        statusColor = color.New(color.FgRed)
    }
    statusColor.Printf("Status: %s (%d)\n", resp.Status, resp.StatusCode)

    // Response time
    color.New(color.FgCyan).Printf("Response Time: %v\n", resp.ResponseTime)

    // Headers
    if len(resp.Headers) > 0 {
        color.New(color.FgHiBlue).Println("Headers:")
        for k, v := range resp.Headers {
            fmt.Printf("  %s: %s\n", k, v)
        }
    }

    // Body: try to pretty print JSON, otherwise print as-is
    if resp.Body != "" {
        fmt.Println()
        color.New(color.FgMagenta).Println("Response Body:")

        var pretty bytesOrString
        if json.Valid([]byte(resp.Body)) {
            var indented bytesOrString
            var raw interface{}
            if err := json.Unmarshal([]byte(resp.Body), &raw); err == nil {
                b, err := json.MarshalIndent(raw, "", "  ")
                if err == nil {
                    pretty = bytesOrString(b)
                } else {
                    pretty = bytesOrString(resp.Body)
                }
            } else {
                pretty = bytesOrString(resp.Body)
            }
            fmt.Println(string(pretty))
        } else {
            fmt.Println(resp.Body)
        }
    }
}

// bytesOrString is a tiny helper type so we can treat both []byte and string easily
type bytesOrString []byte

func (b bytesOrString) String() string { return string(b) }

// SaveRequestHistory appends a request record to a simple history file in the user's home dir.
// This is best-effort and will not fail the command if it cannot write.
func SaveRequestHistory(req APIRequest) error {
    home, err := os.UserHomeDir()
    if err != nil {
        return err
    }
    dir := filepath.Join(home, ".resterx")
    _ = os.MkdirAll(dir, 0o755)
    fp := filepath.Join(dir, "history.ndjson")

    f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
    if err != nil {
        return err
    }
    defer f.Close()

    // enrich with timestamp
    out := map[string]interface{}{
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "request":   req,
    }
    enc, err := json.Marshal(out)
    if err != nil {
        return err
    }
    _, _ = f.Write(append(enc, '\n'))
    return nil
}

// ShowHistory prints the last `limit` history entries (ndjson) from the history file.
func ShowHistory(limit int) error {
    home, err := os.UserHomeDir()
    if err != nil {
        return err
    }
    fp := filepath.Join(home, ".resterx", "history.ndjson")
    data, err := os.ReadFile(fp)
    if err != nil {
        return err
    }
    lines := nonEmptyLines(string(data))
    start := 0
    if len(lines) > limit {
        start = len(lines) - limit
    }
    for i := start; i < len(lines); i++ {
        fmt.Printf("%s\n", lines[i])
    }
    return nil
}

func nonEmptyLines(s string) []string {
    var out []string
    var tmp string
    for _, line := range splitLines(s) {
        if line != "" {
            out = append(out, line)
        } else if tmp != "" {
            out = append(out, tmp)
            tmp = ""
        }
    }
    return out
}

func splitLines(s string) []string {
    // simple split by \n
    var res []string
    curr := ""
    for _, r := range s {
        if r == '\n' {
            res = append(res, curr)
            curr = ""
            continue
        }
        curr += string(r)
    }
    if curr != "" {
        res = append(res, curr)
    }
    return res
}
