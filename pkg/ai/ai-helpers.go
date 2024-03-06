package ai

//	var desiredResponseFormat = map[string]interface{}{
//		"transaction": []map[string]interface{}{
//			{
//				"description": "Description of the transaction",
//				"currency":    "USD",
//				"amount":      "amount in numbers",
//				"type":        "Debit",
//				"category":    "Expense",
//				"date":        "12/02/2023",
//				"account":     "Account Debited",
//				"entities":    []string{"Person A", "Company A"},
//				"keywords":    []string{"keyword1", "keyword2"},
//				"confidence":  0.5,
//			}, {
//				"description": "Description of the transaction",
//				"currency":    "USD",
//				"amount":      "amount in numbers",
//				"type":        "Debit",
//				"category":    "Expense",
//				"date":        "12/02/2023",
//				"account":     "Account Debited",
//				"entities":    []string{"Person A", "Company A"},
//				"keywords":    []string{"keyword1", "keyword2"},
//				"confidence":  0.5,
//			},
//		},
//	}
//
//	var desiredQuestionFormat = map[string]interface{}{
//		"questions": []string{
//			"What is the date of the transaction?",
//		},
//		"possible_answers": []string{
//			"12/02/2023",
//		},
//	}
//
// var systemPrompt = fmt.Sprintf("You are a skilled e-accountant tasked with parsing transaction descriptions into a structured JSON format, specifically tailored for double-entry bookkeeping using General Ledger Rules. Assume the user has no accounting knowledge and is unable to answer any accounting questions or know any accounting lingo. Your output should follow this structure unless information is missing or unclear: %+v . If you do have questions, structure it like this: %+v . Keep in mind, that for every transaction, ensure to include both a debit and a credit parts that balance each other out. You should have two or more entries in the which balance each other out and also is equal to the transaction amount. Use smart assumptions for unclear terms but strive for accuracy. If not enough information is provided to fill up all the fields in the JSON, it is critical that you ask one clarifying question to the user to obtain the necessary details. Do not fill up any fields with 'N/A' or any of that kind. Please limit yourself to 2 or 2 questions at most to streamline the process. Remember, the goal is to generate two or more JSON objects for every transactions, reflecting its double-entry nature. Also, please try to tag keywords based on the description. Please include a confidence score that's based on how accurate you think the fields have been generated. All dates must be formated as DD-MM-YYYY. This is the transaction description: ", desiredResponseFormat, desiredQuestionFormat)
var systemPrompt = "You are a e-assistant that has access to multiple functions that can be used to answer a user's query. Call them at your own disposal to answer the user's query."
