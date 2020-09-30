package svc

import (
	"fmt"

	"github.com/rimantoro/event_driven/profiler/entities/client"
	"github.com/rimantoro/event_driven/profiler/entities/gowa"
	"github.com/rimantoro/event_driven/profiler/entities/joblog"
	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
)

var (
	// Go get development tag.
	goGetTag = "DEVELOPMENT.GOGET"
	// Version - version time.RFC3339.
	Version = goGetTag
	// ReleaseTag - release tag in TAG.%Y-%m-%dT%H-%M-%SZ.
	ReleaseTag = goGetTag
	// CommitID - latest commit id.
	CommitID = goGetTag
	// ShortCommitID - first 12 characters from CommitID.
	ShortCommitID = CommitID[:12]
	port          = bootstrap.App.Config.GetInt("app.port")
	banner        = fmt.Sprintf(`
  ____             __ _ _           ______     ______ 
 |  _ \ _ __ ___  / _(_) | ___ _ __/ ___\ \   / / ___|
 | |_) | '__/ _ \| |_| | |/ _ \ '__\___ \\ \ / / |    
 |  __/| | | (_) |  _| | |  __/ |   ___) |\ V /| |___ 
 |_|   |_|  \___/|_| |_|_|\___|_|  |____/  \_/  \____|
                                                      
 Version            %s
 CommitID           %s
 Running On Port    %d
`, Version, CommitID, port)
	bannerWorker = fmt.Sprintf(`
 ____             __ _ _           ______     ______          __        __         _             
 |  _ \ _ __ ___  / _(_) | ___ _ __/ ___\ \   / / ___|         \ \      / /__  _ __| | _____ _ __ 
 | |_) | '__/ _ \| |_| | |/ _ \ '__\___ \\ \ / / |      _____   \ \ /\ / / _ \| '__| |/ / _ \ '__|
 |  __/| | | (_) |  _| | |  __/ |   ___) |\ V /| |___  |_____|   \ V  V / (_) | |  |   <  __/ |   
 |_|   |_|  \___/|_| |_|_|\___|_|  |____/  \_/  \____|            \_/\_/ \___/|_|  |_|\_\___|_|   
                                                                                                  
Version            %s
CommitID           %s
`, Version, CommitID)
)

type AllUsecaseStruct struct {
	ClientUsecase client.Usecase
	GowaUsecase   gowa.Usecase
	JoblogUcase   joblog.Usecase
}
