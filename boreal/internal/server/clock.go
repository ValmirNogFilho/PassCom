package server

import "github.com/google/uuid"

// CompareClock compara dois relógios de tempo e retorna a relação entre eles.
//
// Retorna:
//   - EQUAL: se os relógios de tempo são iguais
//   - LESS: se o relógio do sistema é mais antigo que o recebido
//   - GREATER: se o relógio do sistema é mais novo que o recebido
//   - CONCURRENT: se os relógios de tempo são concorrentes
func (s *System) CompareClock(firstClock map[uuid.UUID]int, secondClock map[uuid.UUID]int) int {
	// Assuma que VC(x) é o relógio do sistema e VC(y) é o relógio recebido.
	// VC(x) < VC(y) ⇔ ∀z[VC(x)[z] ≤ VC(y)[z]] e ∃w[VC(x)[w] < VC(y)[w]]
	// Lê-se: VC(x) é mais antigo que VC(y) se, para todo z de VC(x), eles são menores ou iguais
	// para o correspondente z em VC(y), e existe um w onde VC(x)[w] é estritamente menor que VC(y)[w].
	// Se a condição acima for atendida, então VC(y) é mais novo que VC(x).
	// Se a condição acima for atendida para VC(y), então VC(y) é mais antigo que VC(x).
	// Senão, se ∃z'[VC(x)[z'] > VC(y)[z']], então VC(x) e VC(y) são concorrentes.
	// Senão, VC(x) e VC(y) são iguais.

	vx := firstClock
	vy := secondClock

	isLess := false
	isGreater := false

	// Itera sobre as chaves em ambos os relógios
	for id, x := range vx {
		y, exists := vy[id]
		if !exists {
			y = 0 // Se o ID não existe em vy, assume que seu valor é 0
		}

		if x < y {
			isLess = true
		} else if x > y {
			isGreater = true
		}

		// Se ambos isLess e isGreater são verdadeiros, eles são concorrentes
		if isLess && isGreater {
			return CONCURRENT
		}
	}

	// Verifica quaisquer IDs em vy que estão ausentes em vx
	for id, y := range vy {
		if _, exists := vx[id]; !exists && y > 0 {
			isLess = true
			if isGreater {
				return CONCURRENT
			}
		}
	}

	// Agora determina a relação com base nas flags
	if isLess {
		return NEWER
	}
	if isGreater {
		return OLDER
	}

	// Se nenhuma flag foi definida, então eles são iguais
	return EQUAL
}

func (s *System) UpdateClock(receivedClock map[uuid.UUID]int) {
	for id, timestamp := range receivedClock {
		if _, exists := s.VectorClock[id]; !exists || timestamp > s.VectorClock[id] {
			s.VectorClock[id] = timestamp
		}
	}
}
