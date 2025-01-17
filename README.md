# Flashcards

## Project description

The flashcards application aims to provide a fun platform to help learners memorize information using interactive cards. It can be used to create, manage and play learning sessions tailored to each user.

## Functional rules
#### Managing flashcards
- Adding a flashcard :
  - A flashcard contains :
  - A question.
  - 4 possible answers
  - A correct answer.
  - Tags to associate the flashcard with one or more categories (e.g. math, history).
#### Managing game sessions
- Launch a game session for a student:
  - Select a flashcard category.
  - Generate a session containing 5 flashcards chosen at random from the selected category.
  - A game session has no time limit.
  - Respond to a list of flashcards:
  - Flashcards must be presented one by one.
 - Each response is recorded in the session as correct or incorrect.

#### Saving data
Record:
  - The student's score.
  - The list of answers provided by the student for each session.
  - Session status:
  - Next question.
  - Session completed or not.




## Models
```
flashcard {
  answer
  responses []response
  numRightResponse
  tags []string
}

responseCard {
  id int (entre 1 et 4)
  proposal string
}
```

```
session {
  studentID
  SessionID
  score
  Category
  flashcardList
  proposalList 
  isFinished
}

SessionState {
  nextCardId
  score
  isFinished
}

QuestionResponseBody {
  flashcardId
  numeroResponse
}

```



## Ressources

#### /flashcards
```
Create
update
search
getbyID
```
#### /sessions
```
Create => studentID & catégorie 
```
#### /sessions/:id/state
```
Get
  => Params : idSession
  => Reponse : idProchaineCarte, score, isFinished
```
#### /sessions/:id/answer
```
Post 
  => Params :idSession
  => Body : idCard, numeroReponse
  => Reponse : idProchaineCarte, score, isFinished
  ```