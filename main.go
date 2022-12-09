package main

import (
	"container/list"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string = ""
)

func init() {
	if Token == "" {
		flag.StringVar(&Token, "token", "", "token")
		flag.Parse()
	}
	if Token == "" {
		Token = os.Getenv("TOKEN")
		if Token == "" {
			panic("You need to input the token.")
		}
	}
}

type MsgHistory struct {
	store    map[string]*list.Element
	list     *list.List
	notifier map[string]chan *discordgo.Message
	sync.RWMutex
}

func NewMsgHistory() *MsgHistory {
	mh := &MsgHistory{}
	mh.list = list.New()
	mh.store = make(map[string]*list.Element, 0)
	mh.notifier = make(map[string]chan *discordgo.Message, 0)
	return mh
}

func (mh *MsgHistory) AddMsg(msg *discordgo.Message) {
	mh.Lock()
	defer mh.Unlock()
	if elem, exists := mh.store[msg.ID]; exists {
		mh.list.MoveToFront(elem)
		return
	} else {
		mh.list.PushFront(msg)
		mh.NotifyAll(msg)
	}
}

func (mh *MsgHistory) GC(maxLifeTime time.Duration) {
	time.AfterFunc(
		maxLifeTime,
		func() {
			mh.Lock()
			defer mh.Unlock()
			elem := mh.list.Back()
			if elem == nil {
				mh.GC(maxLifeTime)
				return
			}
			msg := elem.Value.(*discordgo.Message)
			var lstModTime *time.Time
			if editTime := msg.EditedTimestamp; editTime != nil {
				lstModTime = editTime
			} else {
				lstModTime = &msg.Timestamp
			}
			if time.Now().After(lstModTime.Add(maxLifeTime)) {
				mh.list.Remove(elem)
				delete(mh.store, msg.ID)
				fmt.Printf("A data has been cleared from the 'MsgHistory' cache, id:%s, content:%s\n", msg.ID, msg.Content)
			}
			mh.GC(maxLifeTime)
		},
	)
}

func (mh *MsgHistory) NotifyAll(m *discordgo.Message) {
	go func() {
		mh.RLock()
		defer mh.RUnlock()
		for _, ch := range mh.notifier {
			ch <- m
		}
	}()
}

func (mh *MsgHistory) NewFilter(
	name string,
	criteriaFunc func(m *discordgo.Message) bool,
	max int,
	timeout time.Duration,
	loop bool,
	callbackFunc func([]*discordgo.Message),
) {
	for {
		ch := make(chan *discordgo.Message)

		mh.Lock()
		if _, exists := mh.notifier[name]; exists {
			panic("has existed")
		}
		mh.notifier[name] = ch
		mh.Unlock()
		fmt.Printf("Filter: %q Start\n", name)
		collect := make([]*discordgo.Message, 0)

		if timeout != -1 {
			time.AfterFunc(timeout, func() {
				fmt.Printf("timeout: %q\n", name)
				close(ch)
			})
		}

		for {
			msg, isOpen := <-ch
			isDone := false
			if !isOpen {
				isDone = true
			} else {
				if criteriaFunc(msg) {
					collect = append(collect, msg)
				}
				if max != -1 && len(collect) >= max {
					isDone = true
				}
			}
			if isDone {
				mh.Lock()
				delete(mh.notifier, name)
				mh.Unlock()
				callbackFunc(collect)
				if loop {
					break
				} else {
					return
				}
			}
		}
	}
}

func main() {
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages

	msgHistory := NewMsgHistory()
	go msgHistory.GC(10 * time.Second)
	session.AddHandler(func(s *discordgo.Session, curMsg *discordgo.MessageCreate) {
		appID := s.State.User.ID
		if curMsg.Author.ID == appID {
			return
		}
		msgHistory.AddMsg(curMsg.Message)
	})

	go msgHistory.NewFilter("легион", func(receiveMsg *discordgo.Message) bool {
		return strings.Contains(strings.ToLower(receiveMsg.Content), "легион")
	}, -1, 3*time.Second, true,
		func(collectMsg []*discordgo.Message) {
			if len(collectMsg) == 0 {
				return
			}
			for _, msg := range collectMsg {
				for _, emoji := range []string{"🇱", "🇪", "🇬", "🇮", "🇴", "🇳"} {
					_ = session.MessageReactionAdd(msg.ChannelID, msg.ID, emoji)
					time.Sleep(200 * time.Millisecond)
				}
			}
		})

	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-chanSignal
	_ = session.Close()
}
