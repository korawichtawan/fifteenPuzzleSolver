package fifteenPuzzleSolver

import (
	"container/heap"
	"math"
)

type Position struct {
	x int
	y int
}

type State struct {
	board [4][4]int
	answer []string
	emptyPosition Position
}

type void struct {}

var member void;

var minMove int = math.MaxInt32;
var bestAns []string;
var directions [4]string = [4]string{"t","l","b","r"}


func zeroPosition(board *[4][4]int) (Position) {
	for i,_ := range board {
		for j,val := range board[i] {
			if val == 0 {
				return Position{i,j};
			}
		}
	}
	return Position{-1,-1};
}

func lastMove(s *State) string {
	if len(s.answer) == 0 {
		return ""
	}
	return s.answer[len(s.answer)-1]
}

func canMove(s State, direction string) bool {
	x, y := s.emptyPosition.x,s.emptyPosition.y;

	previousMove := lastMove(&s)

	switch direction {
	case "t":
		if x == 0 || (previousMove == "b") {
			return false
		}
	case "b":
		if x == 3 || (previousMove == "t") {
			return false
		}
	case "r":
		if y == 3 || (previousMove == "l") {
			return false
		}
	case "l":
		if y == 0 || (previousMove == "r") {
			return false
		}
	}

	return true;
}

func swap(board [4][4]int, x1, y1, x2, y2 int) [4][4]int {
	newBoard := board;
	newBoard[x1][y1], newBoard[x2][y2] = newBoard[x2][y2], newBoard[x1][y1];
	return newBoard;
}


func calCost(board *[4][4]int) int{
	var cost float64;
	for i:=0;i<4;i++ {
		for j:=0;j<4;j++{
			var x,y int;
			if board[i][j] == 0 {
				x,y = 3,3
			} else {
				x = (board[i][j] - 1 ) / 4;
				y = (board[i][j] - 1 ) % 4;
			}
			cost += math.Abs(float64(x-i)) + math.Abs(float64(y-j))
		}
	}
	return int(cost);
}

func tryMove(s State, direction string,ch chan State) {
	canMove := canMove(s, direction);
	if canMove == false {
		ch <- State{s.board,[]string{"invalid"},s.emptyPosition}
		return;
	}
	x1, y1 := s.emptyPosition.x,s.emptyPosition.y;
	var x2, y2 int;
	switch direction {
	case "t":
		x2 = x1 - 1;
		y2 = y1;
	case "b":
		x2 = x1 + 1;
		y2 = y1;
	case "r":
		x2 = x1;
		y2 = y1 + 1;
	case "l":
		x2 = x1;
		y2 = y1 - 1;
	}
	newBoard := swap(s.board, x1, y1, x2, y2);
	newAnswer := make([]string,0);
	newAnswer = append(newAnswer, s.answer...)
	newAnswer = append(newAnswer, direction)
	newEmptyPosition := Position{x2, y2};
	ch <- State{newBoard,newAnswer,newEmptyPosition}
	return;
}



func Solve(board [4][4]int) (int,[]string){
	ch := make(chan State,4)
	set := make(map[[4][4]int]void);

	set[board] = member

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	item := &Item{
		value: State{board,[]string{},zeroPosition(&board)},
		priority: calCost(&board),
	}
	heap.Push(&pq, item)

	for pq.Len() > 0 {

		topItem := heap.Pop(&pq).(*Item)
		if (calCost(&(topItem.value.board)) == 0) && (len(topItem.value.answer) < minMove) {
			minMove = len(topItem.value.answer)
			bestAns = topItem.value.answer
			// break
		} 
		if (topItem.priority + len(topItem.value.answer) < minMove) && (len(topItem.value.answer) < 40) {

			for _,direction := range directions {
				go tryMove(topItem.value,direction,ch)
			}

			for i:=0;i<4;i++ {
				newState := <- ch
				newCost := calCost((*[4][4]int)(&newState.board))
				_,exist := set[newState.board]
				if (exist == false) && (newState.answer[0] != "invalid") && (newCost + len(newState.answer) < minMove) {
					newItem := &Item{
						value: newState,
						priority: len(newState.answer) + newCost,
					}
					heap.Push(&pq,newItem)
					set[newState.board] = member;
				}
			}
			
		}
	}

	return minMove,bestAns
}


