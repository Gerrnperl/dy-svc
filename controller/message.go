package controller

import (
	"fmt"
	"main/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatMessageResponse struct {
	Response
	Messages []service.Message `json:"message_list,omitempty"`
}

type ChatListResponse struct {
	Response
	Users []*service.FriendUser `json:"user_list,omitempty"`
}

// POST /douyin/message/action/ - 消息操作
// 登录用户对消息的相关操作，目前只支持消息发送
func MessageAction(c *gin.Context) {
	userId, err := GetUserID(c, "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "用户未登录",
		})
		return
	}
	toUserIdStr := c.Query("to_user_id")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "to_user_id 参数错误",
		})
		return
	}
	if c.Query("action_type") != "1" {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "unsupported action type",
		})
		return
	}
	content := c.Query("content")
	err = service.PostMessage(toUserId, userId, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("发送失败: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "发送成功",
	})
}

// GET /douyin/message/chat/ - 聊天记录
// 当前登录用户和其他指定用户的聊天消息记录
func ChatMessage(c *gin.Context) {
	userId, err := GetUserID(c, "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "用户未登录",
		})
		return
	}
	toUserIdStr := c.Query("to_user_id")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "to_user_id 参数错误",
		})
		return
	}
	preMsgTimeStr := c.Query("pre_msg_time")
	preMsgTime, err := strconv.ParseInt(preMsgTimeStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "pre_msg_time 参数错误",
		})
		return
	}
	messages, err := service.GetMessages(userId, toUserId, preMsgTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("获取聊天记录失败: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, ChatMessageResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "获取聊天记录成功",
		},
		Messages: messages,
	})
}

// /douyin/relation/friend/list/ - 用户好友列表
// 所有关注登录用户的粉丝列表。
func ChatList(c *gin.Context) {
	userId, err := GetUserID(c, "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "用户未登录",
		})
		return
	}
	friends, err := service.GetFriends(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("获取好友列表失败: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, ChatListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "获取好友列表成功",
		},
		Users: friends,
	})
}
