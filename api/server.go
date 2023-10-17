package api

import "net/http"

func NewServer(attorney *AttorneyGeneral) chan error {
	finalize := make(chan error, 2)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./api/static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs)) //

	mux.HandleFunc("/api", attorney.ApiHandler)
	mux.HandleFunc("/", attorney.MainHandler)
	mux.HandleFunc("/boards", attorney.BoardsHandler)
	mux.HandleFunc("/board/", attorney.BoardHandler)
	mux.HandleFunc("/collectives", attorney.CollectivesHandler)
	mux.HandleFunc("/collective/", attorney.CollectiveHandler)
	mux.HandleFunc("/drafts", attorney.DraftsHandler)
	mux.HandleFunc("/draft/", attorney.DraftHandler)
	mux.HandleFunc("/edits/", attorney.EditsHandler)
	mux.HandleFunc("/events", attorney.EventsHandler)
	mux.HandleFunc("/event/", attorney.EventHandler)
	mux.HandleFunc("/members", attorney.MembersHandler)
	mux.HandleFunc("/member/", attorney.MemberHandler)
	mux.HandleFunc("/votes/", attorney.VotesHandler)
	mux.HandleFunc("/requestmembership/", attorney.RequestMemberShipVoteHandler)
	mux.HandleFunc("/newdraft", attorney.NewDraft2Handler)
	mux.HandleFunc("/edit", attorney.NewEditHandler)
	mux.HandleFunc("/editview/", attorney.EditViewHandler)
	mux.HandleFunc("/media/", attorney.MediaHandler)
	mux.HandleFunc("/uploadfile", attorney.UploadHandler)
	mux.HandleFunc("/createboard", attorney.CreateBoardHandler)
	mux.HandleFunc("/votecreateboard/", attorney.VoteCreateBoardHandler)
	mux.HandleFunc("/updateboard/", attorney.UpdateBoardHandler)
	mux.HandleFunc("/voteupdateboard/", attorney.UpdateBoardHandler)
	mux.HandleFunc("/updatecollective/", attorney.UpdateCollectiveHandler)
	mux.HandleFunc("/voteupdatecollective/", attorney.VoteUpdateCollectiveHandler)
	mux.HandleFunc("/updateevent/", attorney.UpdateEventHandler)
	mux.HandleFunc("/votecancelevent/", attorney.VoteCancelEventHandler)
	mux.HandleFunc("/votecreateevent/", attorney.VoteCreateEventHandler)
	mux.HandleFunc("/createevent", attorney.CreateEventHandler)
	mux.HandleFunc("/voteupdateevent/", attorney.VoteUpdateEventHandler)
	mux.HandleFunc("/news", attorney.NewsHandler)
	mux.HandleFunc("/connections/", attorney.ConnectionsHandler)
	mux.HandleFunc("/updates", attorney.UpdatesHandler)
	mux.HandleFunc("/pending", attorney.PendingActionsHandler)
	mux.HandleFunc("/createcollective/", attorney.CreateCollectiveHandler)
	mux.HandleFunc("/mymedia", attorney.MyMediaHandler)
	mux.HandleFunc("/myevents", attorney.MyEventsHandler)
	mux.HandleFunc("/detailedvote/", attorney.DetailedVoteHandler)
	// mux.HandleFunc("/member/votes", attorney.VotesHandler)

	return finalize
}
