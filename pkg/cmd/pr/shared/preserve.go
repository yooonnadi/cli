package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cli/cli/pkg/iostreams"
)

func PreserveInput(io *iostreams.IOStreams, state *IssueMetadataState, createErr *error) func() {
	return func() {
		if !state.IsDirty() {
			return
		}

		if *createErr == nil {
			return
		}

		out := io.ErrOut

		// this extra newline guards against appending to the end of a survey line
		fmt.Fprintln(out)

		data, err := json.Marshal(state)
		if err != nil {
			fmt.Fprintf(out, "failed to save input to file: %s\n", err)
			fmt.Fprintln(out, "would have saved:")
			fmt.Fprintf(out, "%v\n", state)
			return
		}

		random := fmt.Sprintf("%x", time.Now().UnixNano())
		random = random[len(random)-5:]
		dumpFilename := fmt.Sprintf("gh%s.json", random)
		dumpPath := filepath.Join(os.TempDir(), dumpFilename)

		err = ioutil.WriteFile(dumpPath, data, 0660)
		if err != nil {
			fmt.Fprintf(out, "failed to save input to file: %s\n", err)
			fmt.Fprintln(out, "would have saved:")
			fmt.Fprintln(out, string(data))
			return
		}

		cs := io.ColorScheme()

		issueType := "pr"
		if state.Type == IssueMetadata {
			issueType = "issue"
		}

		fmt.Fprintf(out, "%s operation failed. input saved to: %s\n", cs.FailureIcon(), dumpPath)
		fmt.Fprintf(out, "resubmit with: gh %s create -j@%s\n", issueType, dumpPath)

		// some whitespace before the actual error
		fmt.Fprintln(out)
	}
}
