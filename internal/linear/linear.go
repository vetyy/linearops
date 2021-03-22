package linear

const (
	Urgent = "Urgent"
)

type Service struct {
	UnstartedStates []string
	CompletedStates []string
}

func NewService(unstartedStates []string, completedStates []string) *Service {
	return &Service{
		UnstartedStates: unstartedStates,
		CompletedStates: completedStates,
	}
}

func (s *Service) NewAction() *Action {
	return &Action{}
}

type ActionState struct {
	Name string `json:"name"`
}

type ActionData struct {
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	State         *ActionState `json:"state"`
	PriorityLabel string       `json:"priorityLabel"`
}

type Action struct {
	Action string      `json:"action"`
	Data   *ActionData `json:"data"`
	URL    string      `json:"url"`
}

func (s *Service) IsUrgent(action *Action) bool {
	return action.Data.PriorityLabel == Urgent
}

func (s *Service) IsCompleted(action *Action) bool {
	state := action.Data.State.Name
	for _, s := range s.CompletedStates {
		if state == s {
			return true
		}
	}
	return false
	//return state == Resolved || state == Postmortem || state == Rejected
}

func (s *Service) IsUnstarted(action *Action) bool {
	state := action.Data.State.Name
	for _, s := range s.UnstartedStates {
		if state == s {
			return true
		}
	}
	return false
}
