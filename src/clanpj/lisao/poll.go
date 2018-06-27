package main

import (
	"log"
	"sync"
	"time"
)

// PollGithubForever polls github for updates to the bot repo - if any are
// found, a build task for the new commits is pushed.
// TODO(guy) add context here for cancelling
func (s *State) PollGithubForever(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// TODO(guy) this should be persisted in mysql.
	refToCommit := make(map[string]string)

	for {
		refs, err := s.githubClient.GetRefs()
		if err != nil {
			log.Printf("Error getting github refs: %v", err)
			goto Sleep
		}

		for _, ref := range refs {
			oldCommit, ok := refToCommit[ref.Ref]

			// Has the commit changed? I.e. do we need to make a new build.
			// TODO(guy) check DB for whether this commit has been built before.
			if !ok || oldCommit != ref.Object.Sha {
				refToCommit[ref.Ref] = ref.Object.Sha

				// TODO(guy) actually push the build
				log.Printf("PUSHING BUILD FOR %s/%s", ref.Ref, ref.Object.Sha)
			}
		}

	Sleep:
		time.Sleep(time.Second * 30)
	}
}
