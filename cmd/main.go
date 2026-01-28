package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Daemon0x00000000/lql/pkg/lql"
)

// User représente un utilisateur du système
type User struct {
	ID       string
	Name     string
	Email    string
	Age      int
	Role     string
	Status   string
	Tags     []string
	Created  time.Time
	LastSeen time.Time
}

func main() {
	fmt.Println("===============================================================")
	fmt.Println("              LQL - Logical Query Language Demo")
	fmt.Println("===============================================================\n")

	// Mock
	users := []User{
		{
			ID:       "USR001",
			Name:     "Alice Martin",
			Email:    "alice.martin@example.com",
			Age:      28,
			Role:     "admin",
			Status:   "active",
			Tags:     []string{"premium", "verified"},
			Created:  time.Now().AddDate(0, -6, 0),
			LastSeen: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:       "USR002",
			Name:     "Bob Dupont",
			Email:    "bob.dupont@test.fr",
			Age:      35,
			Role:     "moderator",
			Status:   "active",
			Tags:     []string{"verified"},
			Created:  time.Now().AddDate(-1, 0, 0),
			LastSeen: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:       "USR003",
			Name:     "Charlie Smith",
			Email:    "charlie@guest.com",
			Age:      22,
			Role:     "user",
			Status:   "pending",
			Tags:     []string{"new"},
			Created:  time.Now().AddDate(0, 0, -2),
			LastSeen: time.Now().Add(-10 * time.Minute),
		},
		{
			ID:       "USR004",
			Name:     "Diana Lopez",
			Email:    "diana.lopez@example.com",
			Age:      31,
			Role:     "admin",
			Status:   "inactive",
			Tags:     []string{"premium"},
			Created:  time.Now().AddDate(0, -3, 0),
			LastSeen: time.Now().AddDate(0, 0, -10),
		},
		{
			ID:       "USR005",
			Name:     "Eve Johnson",
			Email:    "eve.j@corporate.org",
			Age:      45,
			Role:     "user",
			Status:   "active",
			Tags:     []string{"corporate", "verified", "premium"},
			Created:  time.Now().AddDate(-2, 0, 0),
			LastSeen: time.Now().Add(-30 * time.Minute),
		},
	}

	// ---------------------------------------------------------------
	// Exemple 1 : Filtrage simple
	// ---------------------------------------------------------------
	fmt.Println("\n[Exemple 1] Filtrage simple")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("Query: status=active^age>25")
	runQuery("status=active^age>25", users)

	// ---------------------------------------------------------------
	// Exemple 2 : Utilisation de IN
	// ---------------------------------------------------------------
	fmt.Println("\n[Exemple 2] Opérateur IN")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("Query: roleINadmin,moderator^status=active")
	runQuery("roleINadmin,moderator^status=active", users)

	// ---------------------------------------------------------------
	// Exemple 3 : Recherche par email avec regex
	// ---------------------------------------------------------------
	fmt.Println("\n[Exemple 3] Validation email avec regex")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println(`Query: emailMATCHES'^[a-z]+\.[a-z]+@example\.com$'`)
	runQuery(`emailMATCHES'^[a-z]+\.[a-z]+@example\.com$'`, users)

	// ---------------------------------------------------------------
	// Exemple 4 : Requête complexe avec parenthèses
	// ---------------------------------------------------------------
	fmt.Println("\n[Exemple 4] Requête complexe avec groupement")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("Query: (role=admin^status=active)^OR(role=moderator^age>30)")
	runQuery("(role=admin^status=active)^OR(role=moderator^age>30)", users)

	// ---------------------------------------------------------------
	// Exemple 5 : Recherche de texte
	// ---------------------------------------------------------------
	fmt.Println("\n[Exemple 5] Recherche de texte avec CONTAINS")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("Query: nameCONTAINSMartin^ORnameCONTAINSLopez")
	runQuery("nameCONTAINSMartin^ORnameCONTAINSLopez", users)

	// ---------------------------------------------------------------
	// Exemple 6 : Opérateurs logiques avancés (XOR)
	// ---------------------------------------------------------------
	fmt.Println("\n[Exemple 6] XOR - Exactement une condition vraie")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("Query: role=admin^XORstatus=active")
	runQuery("role=admin^XORstatus=active", users)

	// ---------------------------------------------------------------
	// Exemple 7 : Filtre ultra-complexe
	// ---------------------------------------------------------------
	fmt.Println("\n[Exemple 7] Filtre ultra-complexe")
	fmt.Println("---------------------------------------------------------------")
	complexQuery := "(roleINadmin,moderator^status=active)^OR(age>40^emailCONTAINScorporate)^OR(idSTARTSWITHUSR00^age<30)"
	fmt.Println("Query:", complexQuery)
	fmt.Println("Description: Admins/modos actifs OU utilisateurs 40+ avec email corporate OU jeunes utilisateurs avec ID spécifique")
	runQuery(complexQuery, users)

	// ---------------------------------------------------------------
	// Exemple 8 : Performance test
	// ---------------------------------------------------------------
	fmt.Println("\n[Exemple 8] Test de performance")
	fmt.Println("---------------------------------------------------------------")
	performanceTest()

	// ---------------------------------------------------------------
	// Exemple 9 : Visualisation de l'AST
	// ---------------------------------------------------------------
	fmt.Println("\n[Exemple 9] Visualisation de l'arbre syntaxique (AST)")
	fmt.Println("---------------------------------------------------------------")
	visualizeAST()

	fmt.Println("\n===============================================================")
	fmt.Println("                    Fin de la démo")
	fmt.Println("===============================================================")
}

// runQuery exécute une requête LQL sur une liste d'utilisateurs
func runQuery(query string, users []User) {
	start := time.Now()

	// Parser la requête
	ast, err := lql.Parse(query)
	if err != nil {
		fmt.Printf("[ERREUR] Erreur de parsing: %v\n", err)
		return
	}

	parseTime := time.Since(start)

	// Évaluer pour chaque utilisateur
	matches := []User{}
	evalStart := time.Now()

	for _, user := range users {
		data := userToMap(user)
		if ast.Eval(data) {
			matches = append(matches, user)
		}
	}

	evalTime := time.Since(evalStart)
	totalTime := time.Since(start)

	// Afficher les résultats
	fmt.Printf("\n[TEMPS] Parse: %v | Eval: %v | Total: %v\n", parseTime, evalTime, totalTime)
	fmt.Printf("[RESULTATS] %d utilisateur(s) trouvé(s):\n", len(matches))

	for _, user := range matches {
		fmt.Printf("  - %s (%s) - %s, %d ans - Role: %s - Tags: %v\n",
			user.Name, user.ID, user.Status, user.Age, user.Role, user.Tags)
	}
}

// userToMap convertit un User en map pour LQL
func userToMap(user User) map[lql.Field]interface{} {
	return map[lql.Field]interface{}{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"age":     user.Age,
		"role":    user.Role,
		"status":  user.Status,
		"tags":    strings.Join(user.Tags, ","),
		"created": user.Created.Format("2006-01-02"),
	}
}

// performanceTest teste les performances avec des requêtes répétées
func performanceTest() {
	query := "(roleINadmin,moderator^status=active)^OR(age>30^emailCONTAINSexample)"
	iterations := 10000

	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Iterations: %d\n\n", iterations)

	// Test 1: Parse + Eval à chaque fois
	start := time.Now()
	for i := 0; i < iterations; i++ {
		ast, _ := lql.Parse(query)
		data := map[lql.Field]interface{}{
			"role":   "admin",
			"status": "active",
			"age":    35,
			"email":  "test@example.com",
		}
		ast.Eval(data)
	}
	timeWithParse := time.Since(start)

	// Test 2: Parse une fois, Eval plusieurs fois
	ast, _ := lql.Parse(query)
	start = time.Now()
	for i := 0; i < iterations; i++ {
		data := map[lql.Field]interface{}{
			"role":   "admin",
			"status": "active",
			"age":    35,
			"email":  "test@example.com",
		}
		ast.Eval(data)
	}
	timeEvalOnly := time.Since(start)

	fmt.Printf("Parse + Eval à chaque fois: %v (%.2f μs/op)\n",
		timeWithParse, float64(timeWithParse.Microseconds())/float64(iterations))
	fmt.Printf("Parse une fois + Eval:      %v (%.2f μs/op)\n",
		timeEvalOnly, float64(timeEvalOnly.Microseconds())/float64(iterations))
	fmt.Printf("Gain de performance:        %.1fx plus rapide\n",
		float64(timeWithParse)/float64(timeEvalOnly))
}

// visualizeAST affiche l'arbre syntaxique d'une requête complexe
func visualizeAST() {
	query := "sys_id=123^OR(testINhello,world^XORmeCONTAINS'urgent'^sys_idMATCHES'\\d+')"
	fmt.Printf("Query: %s\n\n", query)

	ast, err := lql.Parse(query)
	if err != nil {
		fmt.Printf("[ERREUR] %v\n", err)
		return
	}

	fmt.Println("Arbre syntaxique (AST):")
	fmt.Println(ast.String())

	fmt.Println("\nÉvaluation avec des données de test:")
	testCases := []struct {
		name string
		data map[lql.Field]interface{}
	}{
		{
			name: "Cas 1: sys_id match",
			data: map[lql.Field]interface{}{
				"sys_id": "123",
				"test":   "other",
				"me":     "normal",
			},
		},
		{
			name: "Cas 2: test IN match",
			data: map[lql.Field]interface{}{
				"sys_id": "456",
				"test":   "hello",
				"me":     "normal",
			},
		},
		{
			name: "Cas 3: me CONTAINS match (XOR avec sys_id MATCHES)",
			data: map[lql.Field]interface{}{
				"sys_id": "abc",
				"test":   "other",
				"me":     "This is urgent!",
			},
		},
		{
			name: "Cas 4: Aucun match",
			data: map[lql.Field]interface{}{
				"sys_id": "abc",
				"test":   "other",
				"me":     "normal",
			},
		},
	}

	for _, tc := range testCases {
		result := ast.Eval(tc.data)
		status := "[FAIL]"
		if result {
			status = "[PASS]"
		}
		fmt.Printf("  %s %s: %v\n", status, tc.name, result)
	}
}
