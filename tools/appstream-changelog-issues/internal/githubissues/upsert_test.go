// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package githubissues

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/google/go-github/v83/github"
)

func TestUpsertIssue_CreateWhenMissing(t *testing.T) {
	t.Parallel()

	fake := &fakeIssuesService{
		listIssues: []*github.Issue{},
		createIssue: &github.Issue{
			Number:  github.Ptr(42),
			HTMLURL: github.Ptr("https://github.com/o/r/issues/42"),
		},
	}
	client := &Client{owner: "o", repo: "r", issues: fake}

	result, err := client.CreateOrUpdateIssue(context.Background(), UpsertInput{
		Title:  "AWS SDK service/appstream feature updates available",
		Body:   "<!-- appstream-changelog-issues -->\ncontent",
		Marker: "<!-- appstream-changelog-issues -->",
		Labels: []string{"dependencies", "appstream-sdk-watch", "dependencies"},
	})
	if err != nil {
		t.Fatalf("UpsertIssue returned error: %v", err)
	}

	if result.Action != UpsertActionCreated {
		t.Fatalf("expected action %q, got %q", UpsertActionCreated, result.Action)
	}
	if result.IssueNumber != 42 {
		t.Fatalf("expected issue number 42, got %d", result.IssueNumber)
	}
	if fake.createReq == nil {
		t.Fatal("expected create request to be sent")
	}

	gotLabels := ptrSlice(fake.createReq.Labels)
	wantLabels := []string{"appstream-sdk-watch", "dependencies"}
	if !slices.Equal(gotLabels, wantLabels) {
		t.Fatalf("expected labels %v, got %v", wantLabels, gotLabels)
	}
}

func TestUpsertIssue_NoopWhenMatching(t *testing.T) {
	t.Parallel()

	fake := &fakeIssuesService{
		listIssues: []*github.Issue{
			{
				Number: github.Ptr(11),
				Title:  github.Ptr("AWS SDK service/appstream feature updates available"),
				Body:   github.Ptr("<!-- appstream-changelog-issues -->\ncontent"),
				Labels: []*github.Label{
					{Name: github.Ptr("dependencies")},
					{Name: github.Ptr("appstream-sdk-watch")},
					{Name: github.Ptr("manual-label")},
				},
			},
		},
	}
	client := &Client{owner: "o", repo: "r", issues: fake}

	result, err := client.CreateOrUpdateIssue(context.Background(), UpsertInput{
		Title:  "AWS SDK service/appstream feature updates available",
		Body:   "<!-- appstream-changelog-issues -->\ncontent",
		Marker: "<!-- appstream-changelog-issues -->",
		Labels: []string{"dependencies", "appstream-sdk-watch"},
	})
	if err != nil {
		t.Fatalf("CreateOrUpdateIssue returned error: %v", err)
	}
	if result.Action != UpsertActionNoop {
		t.Fatalf("expected action %q, got %q", UpsertActionNoop, result.Action)
	}
	if fake.editReq != nil {
		t.Fatal("did not expect edit request for noop case")
	}
}

func TestCreateOrUpdateIssue_UpdateBody(t *testing.T) {
	t.Parallel()

	fake := &fakeIssuesService{
		listIssues: []*github.Issue{
			{
				Number: github.Ptr(12),
				Title:  github.Ptr("AWS SDK service/appstream feature updates available"),
				Body:   github.Ptr("<!-- appstream-changelog-issues -->\nold"),
				Labels: []*github.Label{
					{Name: github.Ptr("dependencies")},
					{Name: github.Ptr("appstream-sdk-watch")},
				},
			},
		},
		editIssue: &github.Issue{Number: github.Ptr(12)},
	}
	client := &Client{owner: "o", repo: "r", issues: fake}

	result, err := client.CreateOrUpdateIssue(context.Background(), UpsertInput{
		Title:  "AWS SDK service/appstream feature updates available",
		Body:   "<!-- appstream-changelog-issues -->\nnew",
		Marker: "<!-- appstream-changelog-issues -->",
		Labels: []string{"dependencies", "appstream-sdk-watch"},
	})
	if err != nil {
		t.Fatalf("CreateOrUpdateIssue returned error: %v", err)
	}
	if result.Action != UpsertActionUpdated {
		t.Fatalf("expected action %q, got %q", UpsertActionUpdated, result.Action)
	}
	if fake.editReq == nil || fake.editReq.Body == nil || *fake.editReq.Body != "<!-- appstream-changelog-issues -->\nnew" {
		t.Fatalf("expected edit request body to be updated, got %#v", fake.editReq)
	}
}

func TestCreateOrUpdateIssue_AddMissingLabels(t *testing.T) {
	t.Parallel()

	fake := &fakeIssuesService{
		listIssues: []*github.Issue{
			{
				Number: github.Ptr(13),
				Title:  github.Ptr("AWS SDK service/appstream feature updates available"),
				Body:   github.Ptr("<!-- appstream-changelog-issues -->\ncontent"),
				Labels: []*github.Label{
					{Name: github.Ptr("dependencies")},
					{Name: github.Ptr("manual-label")},
				},
			},
		},
		editIssue: &github.Issue{Number: github.Ptr(13)},
	}
	client := &Client{owner: "o", repo: "r", issues: fake}

	_, err := client.CreateOrUpdateIssue(context.Background(), UpsertInput{
		Title:  "AWS SDK service/appstream feature updates available",
		Body:   "<!-- appstream-changelog-issues -->\ncontent",
		Marker: "<!-- appstream-changelog-issues -->",
		Labels: []string{"dependencies", "appstream-sdk-watch"},
	})
	if err != nil {
		t.Fatalf("CreateOrUpdateIssue returned error: %v", err)
	}
	if fake.editReq == nil || fake.editReq.Labels == nil {
		t.Fatalf("expected edit request labels to be set, got %#v", fake.editReq)
	}

	got := ptrSlice(fake.editReq.Labels)
	want := []string{"appstream-sdk-watch", "dependencies", "manual-label"}
	if !slices.Equal(got, want) {
		t.Fatalf("expected merged labels %v, got %v", want, got)
	}
}

func TestCreateOrUpdateIssue_ErrorOnMultipleMarkerIssues(t *testing.T) {
	t.Parallel()

	fake := &fakeIssuesService{
		listIssues: []*github.Issue{
			{
				Number: github.Ptr(14),
				Body:   github.Ptr("<!-- appstream-changelog-issues -->\na"),
			},
			{
				Number: github.Ptr(15),
				Body:   github.Ptr("<!-- appstream-changelog-issues -->\nb"),
			},
		},
	}
	client := &Client{owner: "o", repo: "r", issues: fake}

	_, err := client.CreateOrUpdateIssue(context.Background(), UpsertInput{
		Title:  "AWS SDK service/appstream feature updates available",
		Body:   "<!-- appstream-changelog-issues -->\ncontent",
		Marker: "<!-- appstream-changelog-issues -->",
		Labels: []string{"dependencies", "appstream-sdk-watch"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateOrUpdateIssue_ValidatesInput(t *testing.T) {
	t.Parallel()

	client := &Client{owner: "o", repo: "r", issues: &fakeIssuesService{}}
	_, err := client.CreateOrUpdateIssue(context.Background(), UpsertInput{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestCreateOrUpdateIssue_MarkerTakesPrecedenceOverTitle(t *testing.T) {
	t.Parallel()

	fake := &fakeIssuesService{
		listIssues: []*github.Issue{
			{
				Number: github.Ptr(21),
				Title:  github.Ptr("AWS SDK service/appstream feature updates available"),
				Body:   github.Ptr("manual issue without marker"),
			},
			{
				Number: github.Ptr(22),
				Title:  github.Ptr("some other title"),
				Body:   github.Ptr("<!-- appstream-changelog-issues -->\nold automation body"),
			},
		},
		editIssue: &github.Issue{
			Number:  github.Ptr(22),
			HTMLURL: github.Ptr("https://github.com/o/r/issues/22"),
		},
	}
	client := &Client{owner: "o", repo: "r", issues: fake}

	result, err := client.CreateOrUpdateIssue(context.Background(), UpsertInput{
		Title:  "AWS SDK service/appstream feature updates available",
		Body:   "<!-- appstream-changelog-issues -->\nnew automation body",
		Marker: "<!-- appstream-changelog-issues -->",
		Labels: []string{"dependencies", "appstream-sdk-watch"},
	})
	if err != nil {
		t.Fatalf("CreateOrUpdateIssue returned error: %v", err)
	}

	if result.Action != UpsertActionUpdated {
		t.Fatalf("expected action %q, got %q", UpsertActionUpdated, result.Action)
	}
	if fake.editNumber != 22 {
		t.Fatalf("expected marker-matched issue number 22 to be edited, got %d", fake.editNumber)
	}
}

type fakeIssuesService struct {
	listIssues []*github.Issue
	listErr    error

	createIssue *github.Issue
	createErr   error
	createReq   *github.IssueRequest

	editIssue  *github.Issue
	editErr    error
	editReq    *github.IssueRequest
	editNumber int
}

func (f *fakeIssuesService) ListByRepo(_ context.Context, _, _ string, _ *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error) {
	if f.listErr != nil {
		return nil, nil, f.listErr
	}
	return f.listIssues, &github.Response{}, nil
}

func (f *fakeIssuesService) Create(_ context.Context, _, _ string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	f.createReq = issue
	if f.createErr != nil {
		return nil, nil, f.createErr
	}
	if f.createIssue == nil {
		return nil, nil, errors.New("create issue not configured")
	}
	return f.createIssue, &github.Response{}, nil
}

func (f *fakeIssuesService) Edit(_ context.Context, _, _ string, number int, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	f.editNumber = number
	f.editReq = issue
	if f.editErr != nil {
		return nil, nil, f.editErr
	}
	if f.editIssue == nil {
		return nil, nil, errors.New("edit issue not configured")
	}
	return f.editIssue, &github.Response{}, nil
}

func ptrSlice(in *[]string) []string {
	if in == nil {
		return nil
	}
	out := make([]string, len(*in))
	copy(out, *in)
	return out
}
