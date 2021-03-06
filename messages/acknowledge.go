package messages

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/github"
)

//type AcknowledgmentStatus unint32

func CheckAcknowledgement(event github.PullRequestEvent) {
	//Check if PR is in the database
	//Add if not (now or in another function)
	log.Print("MADE IT HERE")
	if checkClosed(event) || checkUnassigned(event) || (event.PullRequest.Assignee == nil && checkOpened(event)) {
		//update event in DB to show the PR is no longer open and no more messages will be accepted for that PR ID
		//don't send a message
	} else if event.PullRequest.Assignee != nil && (checkAssigned(event) || checkReviewed(event) || checkEdited(event) || checkLabeled(event) || checkOpened(event)) {
		//send a message with PR id
		//Format string with PR ID
		log.Printf("Close: %s, Unassigned: %s, Opened: %s, Assigned: %s, Reviewed: %s, Edited: %s, Labeled: %s", checkClosed(event), checkUnassigned(event), checkOpened(event),
			checkAssigned(event), checkReviewed(event), checkEdited(event), checkLabeled(event))
		log.Print("MADE IT HERE")
		message := fmt.Sprintf("PR id, %d, URL, %s, Assignee, %s", *event.PullRequest.ID, *event.PullRequest.URL, *event.PullRequest.Assignee.Login)
		log.Print(message)
		err := SendToQueue(message)
		log.Print("SENT TO QUEUE")
		if err != nil {
			log.Printf("Message for event %d not delivered", *event.PullRequest.ID)
		}
	}
	//error
}

func CheckAcknowledgementComment(event github.IssueCommentEvent) {
	if checkCommented(event) {
		message := fmt.Sprintf("PR id, %d, URL, %s, Assignee, %s", *event.Issue.ID, *event.Issue.URL, *event.Issue.Assignee.Login)
		log.Print(message)
		err := SendToQueue(message)
		log.Print("SENT TO QUEUE")
		if err != nil {
			log.Printf("Message for event %d not delivered", *event.Issue.ID)
		}
	}
}

func checkCommented(event github.IssueCommentEvent) bool {
	//check that the issue is not nil and that the issue id is a pr id in the db
	if event.Issue != nil && event.Issue.Assignee != nil {
		return true
	}
	return false
}

func checkAssigned(event github.PullRequestEvent) bool {
	if strings.Compare(*event.Action, "assigned") == 0 {
		//Update PR in DB to accept messages
		return true
	}
	return false
}

func checkUnassigned(event github.PullRequestEvent) bool {
	if strings.Compare(*event.Action, "unassigned") == 0 {
		//Update PR in DB to no longer accept messages until assigned
		return true
	}
	return false
}

func checkReviewed(event github.PullRequestEvent) bool {
	if strings.Compare(*event.Action, "review_requested") == 0 {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			if strings.Compare(*event.PullRequest.Assignee.Name, *event.Sender.Name) == 0 {
				//Update DB
				return true
			}
		}
	}
	return false
}

func checkLabeled(event github.PullRequestEvent) bool {
	if strings.Compare(*event.Action, "labeled") == 0 {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			//if strings.Compare(*event.PullRequest.Assignee.Name, *event.Sender.Name) == 0 {
			//Update DB
			return true
			//}
		}
	}
	return false
}

func checkClosed(event github.PullRequestEvent) bool {
	if strings.Compare(*event.Action, "closed") == 0 {
		//Update DB to not accept messages
		return true
	}
	return false
}

func checkOpened(event github.PullRequestEvent) bool {
	if strings.Compare(*event.Action, "opened") == 0 {
		//check if PR ID is in the DB
		//Create new entry with assignee if not
		if event.PullRequest.Assignee == nil {
			//PR id cannot accept messages (not assigned yet)
		} else {
			//update DB for PR to accept messages
			//Case when PR closed with an assignee but then reopenned
		}
		return true
	}
	return false
}

func checkEdited(event github.PullRequestEvent) bool {
	if strings.Compare(*event.Action, "edited") == 0 {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			//if strings.Compare(*event.PullRequest.Assignee.Name, *event.Sender.) == 0 {
			//Update DB to accept messages
			return true
			//}
		}
	}
	return false
}
