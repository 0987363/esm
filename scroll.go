/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"gopkg.in/cheggaaa/pb.v1"
	"encoding/json"
	log "github.com/cihub/seelog"
)


// Stream from source es instance. "done" is an indicator that the stream is
// over
func (s *Scroll) ProcessScrollResult(c *Config, bar *pb.ProgressBar){

	//update progress bar
	bar.Add(len(s.Hits.Docs))

	// show any failures
	for _, failure := range s.Shards.Failures {
		reason, _ := json.Marshal(failure.Reason)
		log.Errorf(string(reason))
	}

	// write all the docs into a channel
	for _, docI := range s.Hits.Docs {
		c.DocChan <- docI.(map[string]interface{})

		//write to file channel
		if(len(c.DumpInputFile)>0){
			c.FileChan <- docI.(map[string]interface{})
		}
	}
}

func (s *Scroll) Next(c *Config, bar *pb.ProgressBar) (done bool) {

	scroll,err:=c.SrcESAPI.NextScroll(c.ScrollTime,s.ScrollId)
	if err != nil {
		log.Error(err)
		return false
	}

	if scroll.Hits.Docs == nil || len(scroll.Hits.Docs) <= 0 {
		log.Debug("scroll result is empty")
		return true
	}

	scroll.ProcessScrollResult(c,bar)

	//update scrollId
	s.ScrollId=scroll.ScrollId

	return
}

