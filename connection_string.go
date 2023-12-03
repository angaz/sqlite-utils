package sqlu

import (
	"fmt"
	"net/url"
)

type Pragma string
type Param string

var (
	ParamCachePrivate           Param  = "cache=private"
	ParamCacheShared            Param  = "cache=shared"
	ParamImmutable              Param  = "immutable=1"
	ParamModeMemory             Param  = "memory"
	ParamModeRO                 Param  = "ro"
	ParamModeRW                 Param  = "rw"
	ParamModeRWC                Param  = "rwc"
	ParamNoLock                 Param  = "nolock=1"
	ParamPowerSafeOverWriteOff  Param  = "psow=0"
	ParamPowerSafeOverWriteOn   Param  = "psow=1"
	PragmaAutoVacuumFull        Pragma = "auto_vacuum=FULL"
	PragmaAutoVacuumIncremental Pragma = "auto_vacuum=INCREMENTAL"
	PragmaAutoVacuumNone        Pragma = "auto_vacuum=NONE"
	PragmaJournalModeWAL        Pragma = "journal_mode=wal"
	PragmaSynchronousNormal     Pragma = "synchronous=normal"
)

func ParamModeOfFile(filename string) Param {
	return Param("modeof=" + filename)
}

func PragmaBusyTimeout(t int64) Pragma {
	return Pragma(fmt.Sprintf("busy_timeout=%d", t))
}

// PragmaJournalSizeLimit Sets the maximum size of the journal in BYTES.
// Must be set on every attached database separately.
//
// https://www.sqlite.org/pragma.html#pragma_journal_size_limit
func PragmaJournalSizeLimit(n int64) Pragma {
	return Pragma(fmt.Sprintf("journal_size_limit=%d", n))
}

type ConnParams struct {
	Filename string
	Mode     Param
	Pragma   []Pragma
	Attach   []AttachParams
}

type AttachParams struct {
	Filename string
	Database string
	Mode     Param
	Pragma   []Pragma
}

func (p ConnParams) ConnectionString() string {
	query := url.Values{}

	if p.Mode != "" {
		query.Set("mode", string(p.Mode))
	}

	for _, pragma := range p.Pragma {
		query.Add("_pragma", string(pragma))
	}

	for _, attach := range p.Attach {
		query.Add("_attach", attach.ConnectionString())
	}

	return p.Filename + "?" + query.Encode()
}

func (p AttachParams) ConnectionString() string {
	query := url.Values{
		"_name": []string{p.Database},
	}

	if p.Mode != "" {
		query.Set("mode", string(p.Mode))
	}

	for _, pragma := range p.Pragma {
		query.Add("_pragma", string(pragma))
	}

	return p.Filename + "?" + query.Encode()
}
