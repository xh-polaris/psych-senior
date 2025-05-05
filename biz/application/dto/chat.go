package dto

type (
	// ChatStartReq 开始对话请求
	ChatStartReq struct {
		// 开始的时间戳
		Timestamp int64 `json:"timestamp"`
		// 使用者标记
		From string `json:"from"`
		// 语言
		Lang string `json:"lang"`
	}

	// ChatReq 对话请求
	ChatReq struct {
		// 命令, 0对话, -1结束
		Cmd int64  `json:"cmd"`
		Msg string `json:"msg"`
	}

	// ChatEndResp 对话结束响应
	ChatEndResp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	// ChatData 一次流式响应
	ChatData struct {
		Id        uint64 `json:"id"`
		Content   string `json:"content"`
		SessionId string `json:"session_id"`
		Timestamp int64  `json:"timestamp"`
		Finish    string `json:"finish"`
	}

	// ChatHistory 对话记录
	ChatHistory struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	Report struct {
		BasicInfo        BasicInfo        `json:"basic_info" bson:"basic_info"`
		OverviewSummary  OverviewSummary  `json:"overview_summary" bson:"overview_summary"`
		DetailedAnalysis DetailedAnalysis `json:"detailed_analysis" bson:"detailed_analysis"`
	}

	BasicInfo struct {
		AnalysisDate string `json:"analysis_date" bson:"analysis_date"`
	}

	OverviewSummary struct {
		EmotionTone          string   `json:"emotion_tone" bson:"emotion_tone"`
		MainEmotions         []string `json:"main_emotions" bson:"main_emotions"`
		HealthConcerns       string   `json:"health_concerns" bson:"health_concerns"`
		LonelinessLevel      string   `json:"loneliness_level" bson:"loneliness_level"`
		PositiveLifeSigns    string   `json:"positive_life_signs" bson:"positive_life_signs"`
		SocialStatus         string   `json:"social_status" bson:"social_status"`
		AttentionSuggestions string   `json:"attention_suggestions" bson:"attention_suggestions"`
	}

	DetailedAnalysis struct {
		EmotionStatus              EmotionStatus              `json:"emotion_status" bson:"emotion_status"`
		HealthFocus                HealthFocus                `json:"health_focus" bson:"health_focus"`
		PsychologicalSignals       PsychologicalSignals       `json:"psychological_signals" bson:"psychological_signals"`
		InterestAndPositiveLife    InterestAndPositiveLife    `json:"interest_and_positive_life" bson:"interest_and_positive_life"`
		SocialRelationshipStatus   SocialRelationshipStatus   `json:"social_relationship_status" bson:"social_relationship_status"`
		SupportAndAttentionSummary SupportAndAttentionSummary `json:"support_and_attention_summary" bson:"support_and_attention_summary"`
	}

	EmotionStatus struct {
		OverallEmotionTone string   `json:"overall_emotion_tone" bson:"overall_emotion_tone"`
		MainEmotions       []string `json:"main_emotions" bson:"main_emotions"`
		EmotionVariation   string   `json:"emotion_variation" bson:"emotion_variation"`
	}

	HealthFocus struct {
		MentionedHealthIssues     []string        `json:"mentioned_health_issues" bson:"mentioned_health_issues"`
		HealthRiskAlert           HealthRiskAlert `json:"health_risk_alert" bson:"health_risk_alert"`
		MedicalConsultationIntent bool            `json:"medical_consultation_intent" bson:"medical_consultation_intent"`
	}

	HealthRiskAlert struct {
		Exists  bool   `json:"exists" bson:"exists"`
		Details string `json:"details" bson:"details"`
	}

	PsychologicalSignals struct {
		LonelinessDetected       DetectionStatus `json:"loneliness_detected" bson:"loneliness_detected"`
		CognitiveSignalsDetected DetectionStatus `json:"cognitive_signals_detected" bson:"cognitive_signals_detected"`
		MajorLifeEventsMentioned DetectionStatus `json:"major_life_events_mentioned" bson:"major_life_events_mentioned"`
	}

	DetectionStatus struct {
		Exists  bool   `json:"exists" bson:"exists"`
		Example string `json:"example" bson:"example"`
		Details string `json:"details" bson:"details"`
	}

	InterestAndPositiveLife struct {
		Hobbies                  []string         `json:"hobbies" bson:"hobbies"`
		PositiveAttitude         AttitudeStatus   `json:"positive_attitude" bson:"positive_attitude"`
		InitiativeToTryNewThings InitiativeStatus `json:"initiative_to_try_new_things" bson:"initiative_to_try_new_things"`
	}

	AttitudeStatus struct {
		Exists bool   `json:"exists" bson:"exists"`
		Reason string `json:"reason" bson:"reason"`
	}

	InitiativeStatus struct {
		Exists bool `json:"exists" bson:"exists"`
	}

	SocialRelationshipStatus struct {
		FamilyContactFrequency string         `json:"family_contact_frequency" bson:"family_contact_frequency"`
		FriendInteraction      string         `json:"friend_interaction" bson:"friend_interaction"`
		SocialAttitude         SocialAttitude `json:"social_attitude" bson:"social_attitude"`
	}

	SocialAttitude struct {
		Type    string `json:"type" bson:"type"`
		Example string `json:"example" bson:"example"`
	}

	SupportAndAttentionSummary struct {
		EmotionalSupportNeeds  []string `json:"emotional_support_needs" bson:"emotional_support_needs"`
		HealthSupportNeeds     []string `json:"health_support_needs" bson:"health_support_needs"`
		SpecialAttentionPoints []string `json:"special_attention_points" bson:"special_attention_points"`
	}
)
