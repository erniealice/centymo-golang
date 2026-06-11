package centymo

import shared "github.com/erniealice/centymo-golang/domain/shared"

// DataSource is re-exported from domain/shared (centymo restructure). The
// canonical definition lives in the leaf package so entity view packages can
// import it without an entity -> root -> facade import cycle. External root
// consumers (centymo.DataSource) keep resolving via this alias unchanged.
type DataSource = shared.DataSource
