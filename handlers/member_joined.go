package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/halushko/core-go/nats"
	"github.com/halushko/member-handler-go/database"
	"github.com/halushko/tg-bot-go/bot"
)

func StartMemberJoinedListener() {
	processor := func(data []byte) {
		chatId, args, err := nats.ParseTgBotCommand(data)
		if err != nil {
			log.Printf("[ERROR] Помилка при парсингу повідомлення: %v", err)
			return
		}
		if args == nil || len(args) != 2 {
			log.Printf("[ERROR] Має бути 2 аргументи - chatId та userLogin: args=%v", args)
			return
		}

		if chatId != 0 && args[0] != "" {
			userId, errInt := strconv.ParseInt(args[0], 10, 64)
			if errInt != nil {
				log.Printf("[ERROR] Помилка: ID користувача не число: %s", args[0])
				return
			}
			process(chatId, userId, args[1])
		} else {
			log.Printf("[ERROR] Помилка: ID користувача чи ID чату порожні")
		}
	}

	listener := &nats.ListenerHandler{
		Function: processor,
	}

	nats.StartNatsListener(bot.TelegramMemberJoinedQueue, listener)
}

func process(chatId int64, userId int64, userLogin string) {
	log.Printf("[DEBUG] user id=%d userLogin=%s joined to the chat chatId=%d ", userId, userLogin, chatId)

	student := updateUserInfo(userId, userLogin)
	course := updateCourse(chatId)
	updateStudentCourse(student[database.CUserId], course[database.CCourseId])
}

func updateUserInfo(userId int64, userLogin string) map[string]any {
	usr := make(map[string]any)
	usr[database.CUserId] = userId
	if userLogin != "" {
		usr[database.CLogin] = userLogin
	}

	userInfo, sErr := database.DB.SelectByAny(database.TStudents, usr)
	if sErr != nil {
		log.Printf("[ERROR] can't select student id=%d, login=%s: %v", userId, userLogin, sErr)
		return nil
	}

	if userInfo == nil || len(userInfo) == 0 {
		iErr := database.DB.Insert(database.TStudents, usr)
		if iErr != nil {
			log.Printf("[ERROR] can't insert student id=%d, login=%s: %v", userId, userLogin, iErr)
			return nil
		}
	} else if len(userInfo) == 1 {
		usr = setNonempty(usr, userInfo[0])
		uErr := database.DB.UpdateByAny(database.TStudents, usr, usr)
		if uErr != nil {
			log.Printf("[ERROR] can't update student id=%d, login=%s: %v", userId, userLogin, uErr)
			return nil
		}
	} else {
		tryToFixUser()
		log.Printf("[ERROR] duplicate of student %d (%s)", userId, userLogin)
		return nil
	}
	return usr
}

func updateCourse(chatId any) map[string]any {
	crs := make(map[string]any)
	crs[database.CChatId] = chatId

	courses, err := database.DB.SelectByAll(database.TCourses, crs)
	if err != nil {
		log.Printf("[ERROR] can't get course for chat %s: %v", chatId, err)
		return make(map[string]any)
	}

	if courses == nil || len(courses[0]) == 0 {
		log.Printf("[ERROR] can't course for chat %s doesn't exists", chatId)
		return make(map[string]any)
	}

	if len(courses) > 1 {
		coursesInfo := make([]string, 0)
		for _, course := range courses {
			coursesInfo = append(coursesInfo, fmt.Sprintf("%v:%v ", course[database.CCourseId], course[database.CCourseFullName]))
		}
		log.Printf("[ERROR] There are dupbicates for course for chat %s: %v", chatId, coursesInfo)
		return make(map[string]any)
	}

	return courses[0]
}

func tryToFixUser() {

}

func updateStudentCourse(userId any, courseId any) {
	uc := make(map[string]any)
	uc[database.CUserId] = userId
	uc[database.CCourseId] = courseId

	err := database.DB.InsertIfNotExists(database.TStudentsInCourses, uc)
	if err != nil {
		log.Printf("[ERROR] can't insert student %v in course %v: %v", userId, courseId, err)
	}
}

func setNonempty(tg, db map[string]any) map[string]any {
	keys := make(map[string]any)
	res := make(map[string]any)
	for k := range tg {
		keys[k] = nil
	}
	for k := range db {
		keys[k] = nil
	}
	for k := range keys {
		first, second := selectNonempty(tg[k], db[k])
		if second == nil {
			res[k] = first
		} else {
			log.Printf("[INFO] Несинхронызованы даны в ТГ та БД: %v та %v", first, second)
		}
	}
	return res
}

func selectNonempty(tg, db any) (any, any) {
	inTg := tg == nil
	inDb := db == nil

	switch {
	case !inTg && !inDb:
		return nil, nil
	case inTg && inDb:
		return tg, db
	case !inTg:
		return db, nil
	default:
		return tg, nil
	}
}
