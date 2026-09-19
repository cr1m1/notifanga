package main

import (
	"strings"
	"testing"
)

type subKey struct{ userID, mangaID int }

// fakeRepo is an in-memory NotifangaRepository for service tests.
type fakeRepo struct {
	mangas   map[int]*Manga
	subs     map[subKey]string // last_read_chapter per subscription
	subOrder []subKey
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{mangas: map[int]*Manga{}, subs: map[subKey]string{}}
}

func (f *fakeRepo) subscribe(u *User, m *Manga, lastRead string) {
	f.mangas[m.ID] = m
	k := subKey{u.ID, m.ID}
	if _, ok := f.subs[k]; !ok {
		f.subOrder = append(f.subOrder, k)
	}
	f.subs[k] = lastRead
}

func (f *fakeRepo) UserCreate(u *User) (*User, error)    { return u, nil }
func (f *fakeRepo) UserList() ([]*User, error)           { return nil, nil }
func (f *fakeRepo) MangaCreate(m *Manga) (*Manga, error) { return m, nil }
func (f *fakeRepo) UserListByManga(Manga) ([]*User, error) {
	return nil, nil
}
func (f *fakeRepo) UpdateManga(Manga) error                 { return nil }
func (f *fakeRepo) AddMangaToUser(*Manga, *User) error      { return nil }
func (f *fakeRepo) DeleteMangaFromUser(*Manga, *User) error { return nil }

func (f *fakeRepo) MangaList(u User) ([]*Manga, error) {
	var out []*Manga
	for _, k := range f.subOrder {
		if k.userID == u.ID {
			out = append(out, f.mangas[k.mangaID])
		}
	}
	return out, nil
}

func (f *fakeRepo) SubscriptionList(u User) ([]*Subscription, error) {
	var out []*Subscription
	for _, k := range f.subOrder {
		if k.userID == u.ID {
			out = append(out, &Subscription{Manga: *f.mangas[k.mangaID], LastReadChapter: f.subs[k]})
		}
	}
	return out, nil
}

func (f *fakeRepo) MarkRead(m *Manga, u *User, chapter string) error {
	k := subKey{u.ID, m.ID}
	if _, ok := f.subs[k]; !ok {
		return ErrNotValidChapter
	}
	f.subs[k] = chapter
	return nil
}

func TestSubscriptionUnread(t *testing.T) {
	cases := []struct {
		name     string
		last     string
		lastRead string
		want     bool
	}{
		{"never crawled", "", "", false},
		{"never read", "12", "", true},
		{"behind", "12", "11", true},
		{"caught up", "12", "12", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Subscription{Manga: Manga{LastChapter: c.last}, LastReadChapter: c.lastRead}
			if got := s.Unread(); got != c.want {
				t.Fatalf("Unread() = %v, want %v", got, c.want)
			}
		})
	}
}

func seed(t *testing.T) (*NotifangaService, *fakeRepo, *User, *User) {
	t.Helper()
	repo := newFakeRepo()
	alice := &User{ID: 1, TelegramUserID: 100}
	bob := &User{ID: 2, TelegramUserID: 200}
	repo.subscribe(alice, &Manga{ID: 10, Name: "A", LastChapter: "5"}, "5")
	repo.subscribe(alice, &Manga{ID: 20, Name: "B", LastChapter: "7"}, "6")
	repo.subscribe(alice, &Manga{ID: 30, Name: "C", LastChapter: "1"}, "")
	repo.subscribe(alice, &Manga{ID: 40, Name: "D", LastChapter: ""}, "")
	repo.subscribe(bob, &Manga{ID: 20, Name: "B", LastChapter: "7"}, "")
	return NewNotifangaService(repo), repo, alice, bob
}

func TestListUnread(t *testing.T) {
	svc, _, alice, _ := seed(t)
	unread, err := svc.ListUnread(*alice)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, s := range unread {
		names = append(names, s.Manga.Name)
	}
	if got := strings.Join(names, ","); got != "B,C" {
		t.Fatalf("ListUnread = %q, want %q", got, "B,C")
	}
}

func TestMarkRead(t *testing.T) {
	svc, repo, alice, bob := seed(t)
	if err := svc.MarkRead(&Manga{ID: 20, LastChapter: "7"}, alice); err != nil {
		t.Fatal(err)
	}
	if got := repo.subs[subKey{alice.ID, 20}]; got != "7" {
		t.Fatalf("alice last_read_chapter = %q, want %q", got, "7")
	}
	if got := repo.subs[subKey{bob.ID, 20}]; got != "" {
		t.Fatalf("bob's subscription must be untouched, got %q", got)
	}
	unread, _ := svc.ListUnread(*alice)
	if len(unread) != 1 || unread[0].Manga.Name != "C" {
		t.Fatalf("expected only C unread after marking B, got %d", len(unread))
	}
}

func TestMarkAllRead(t *testing.T) {
	svc, repo, alice, bob := seed(t)
	if err := svc.MarkAllRead(alice); err != nil {
		t.Fatal(err)
	}
	unread, _ := svc.ListUnread(*alice)
	if len(unread) != 0 {
		t.Fatalf("expected no unread, got %d", len(unread))
	}
	if got := repo.subs[subKey{alice.ID, 40}]; got != "" {
		t.Fatalf("manga without chapters must stay '', got %q", got)
	}
	if got := repo.subs[subKey{bob.ID, 20}]; got != "" {
		t.Fatalf("bob's subscription must be untouched, got %q", got)
	}
}

func TestPickByIndices(t *testing.T) {
	list := []*Subscription{{Manga: Manga{Name: "A"}}, {Manga: Manga{Name: "B"}}}
	got := joinNames(pickByIndices(list, []string{"1", "x", "-1", "5", "0"}))
	if got != "B, A" {
		t.Fatalf("pickByIndices = %q, want %q", got, "B, A")
	}
}

func TestNewChapterMessageHintsRead(t *testing.T) {
	msg := NewChapterMessage(Manga{Name: "A", LastChapter: "9", LastChapterUrl: "https://x/9"})
	for _, want := range []string{"9", "A", "https://x/9", "/read"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("message %q lacks %q", msg, want)
		}
	}
}
