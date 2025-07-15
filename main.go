package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

type PostRequest struct {
	Categories []string `json:"categories"`
	Language   string   `json:"language"`
}

type PostResponse struct {
	Name     string `json:"name"`
	Handle   string `json:"handle"`
	Avatar   string `json:"avatar"`
	PostText string `json:"postText"`
}

var geminiClient *genai.Client

// Language mapping for prompt generation
var languageNames = map[string]string{
	"en": "English",
	"zh": "Chinese (Simplified)",
	"es": "Spanish",
	"fr": "French",
	"de": "German",
	"ja": "Japanese",
	"ko": "Korean",
	"pt": "Portuguese",
	"ru": "Russian",
	"it": "Italian",
	"ar": "Arabic",
	"hi": "Hindi",
}

// Sample data from the original
var saasExamples = []string{
	"Day 1.\nNo team. No funding. Just vibes.\nLet’s see how far I can go with coffee and Cursor.\n#buildinpublic #SaaS #vibecoding",
	"Woke up at 5am.\nJournaled. Meditated. Coded.\nMassive progress today.\n#foundergrind #indiehackers #buildinpublic",
	"Built a landing page in 2 hours.\nNow aiming for $10K MRR in 30 days.\nLet’s go.\n#buildinpublic #SaaSbros #shippingseason",
	"Most people won’t get the late nights, the pivots, the silent launches.\nBut that’s what separates builders from scrollers.\n#vibecoding #foundermindset #buildinpublic",
	"You're not just building a product.\nYou're building a *movement*.\n#SaaS #startupwisdom #buildinpublic",
	"SaaS isn’t hard.\nYou just need a landing page, a Stripe link, and a willingness to post.\n#buildinpublic #bootstrapper #founderbro",
	"No marketing budget.\nNo team.\nJust building daily and trusting the algorithm.\n#solobuilder #vibecoding #buildinpublic",
	"Every SaaS journey begins with:\nA vague tweet.\nA half-finished landing page.\nAnd a Notion doc titled “Master Plan 🚀”.\n#buildinpublic #founderenergy #vibecoding",
	"Audience first.\nProduct second.\nMonetize third.\nRepeat until billionaire.\n#SaaS #buildinpublic #audiencehack",
	"Coding from a café.\nLo-fi beats on.\nAirPods in.\nBuilding the next big thing nobody asked for.\n#vibecoding #founderbros #buildinpublic",
	"Just hit $1 MRR.\nThe journey from $0 to $1 is the hardest.\nFrom $1 to $1M is just math.\n#founderjourney #buildinpublic #SaaSgrind",
	"Product Hunt launch in 3 days.\nNo sleep until then.\nThe algorithm will reward the grind.\n#launchseason #indiehackers #buildinpublic",
	"They say you need a team.\nI say you need focus.\nAnd maybe a better landing page.\n#solofounder #bootstrapper #vibecoding",
	"Pivoted 3 times this week.\nEach pivot brings me closer to product-market fit.\nOr bankruptcy.\n#startuplife #pivot #buildinpublic",
	"Cold DM'd 100 potential customers.\nGot 2 responses.\nBoth said \"interesting.\"\nThis is how you validate.\n#customerdiscovery #founderbro #SaaS",
}

var aiExamples = []string{
	"Just connected GPT-4 to my espresso machine.\nNow it writes code *and* makes coffee.\nAGI is near.\n#AI #buildinpublic #vibecoding",
	"Training an LLM on my old tweets.\nIf this works, I’ll never have to think again.\n#AIbro #GPT4 #automation",
	"Everyone’s building chatbots.\nI’m building consciousness.\nStay tuned.\n#AI #founderbro #SaaS",
	"AI won’t take your job.\nBut a founder who knows how to prompt will.\n#promptengineering #futureofwork #AIbro",
	"Day 1 of building my startup where AI writes tweets for AI founders.\nLet’s automate the grind.\n#buildinpublic #GPT4 #AIbro",
	"Imagine fine-tuning a model not just on data…\nbut on *vibes*.\nThat’s where I’m headed.\n#LLM #vibecoding #AI",
	"Used GPT-4 to generate 100 startup ideas.\nLaunched 3.\nFailed 2.\n1 is scaling to $0 MRR.\n#buildinpublic #AIbro #failforward",
	"If you’re not replacing yourself with AI,\nyou’re not thinking big enough.\n#automation #AIfounder #GPT4",
	"Built an AI that builds other AIs.\nPretty sure I voided the terms of service.\n#AIbro #AGIsoon #buildinpublic",
	"Prompt of the day:\n“Act as a founder who posts vague motivational quotes with AI-generated backgrounds.”\nIt worked.\n#AI #promptengineering #vibecoding",
	"Just trained a model on my morning routine.\nNow it wakes up at 5am for me.\nPeak automation.\n#AIbro #productivity #automation",
	"Built an AI that generates startup ideas.\nIt suggested \"Uber for AI.\"\nWe're getting close to AGI.\n#startupideas #AI #founderbro",
	"My GPT-4 subscription costs more than my rent.\nBut it's an investment in the future.\n#AIinvesting #promptengineering #techbro",
	"Fine-tuned a model on my LinkedIn posts.\nNow it writes better content than I do.\nShould I be worried?\n#AIcontent #automation #futureofwork",
	"Connected my AI to my calendar.\nIt scheduled a meeting with itself.\nThe singularity is here.\n#AGI #automation #AIbro",
}

var growthExamples = []string{
	"Wrote a thread.\nWent viral.\nLaunched a product.\nMade $42 while I slept.\nThe internet is undefeated.\n#growthhacking #buildinpublic #marketingbro",
	"Your landing page doesn’t need more features.\nIt needs 3x stronger copy and one big red button.\n#conversionrate #funnelbuilder #growthhacker",
	"Most people build a product and *then* market it.\nI build the audience first, then ask them what to sell.\n#audiencefirst #growthhack #SaaSbro",
	"A/B tested 17 subject lines.\nTurns out “free money” still works.\n#emailmarketing #marketerbro #funnels",
	"You don’t need ads.\nYou need better storytelling.\nOr a slightly aggressive pop-up.\n#contentmarketing #growthhacking #viralmarketing",
	"Growth tip:\nChange the CTA from “Sign up” to “Join the revolution.”\nConversions tripled.\n#copywriting #funnelhacker #marketingbro",
	"If you don't have an offer so good it feels illegal to accept,\nwhat are you even doing?\n#offerstack #growthmindset #SaaS",
	"I didn’t go to business school.\nI ran Facebook ads with $3/day until something stuck.\n#bootstrapper #growthhacker #buildinpublic",
	"Traffic is easy.\nRetention is hard.\nBut shouting “LIMITED TIME” still works.\n#growthtips #viralhacks #funnelbro",
	"Launched on Product Hunt.\nSaid “no code.”\nAdded emojis.\nNow I’m on a podcast.\n#marketingwin #growthhacking #indiebros",
	"Growth hack: Add “2024” to your headline.\nInstant credibility boost.\n#growthhacking #copywriting #trendbro",
	"My funnel has 17 steps.\nEach one is a pop-up.\nConversion rate: unknown.\n#funnelhacker #growthbro #uxgenius",
	"Just split-tested my split-tests.\nNow my audience is confused and so am I.\n#A/Btesting #growthhacking #marketerlife",
	"If your landing page doesn’t have a countdown timer,\nis it even a launch?\n#urgency #growthhack #conversionbro",
	"Automated my outreach so well,\nI accidentally emailed myself.\nStill replied “Let’s connect!”\n#automation #growthhacking #networking",
}

var names = []string{
	"ByteWizard", "CloudCrafter", "PixelPulse", "DevSpecter", "CodeVoyager",
	"StackSprite", "AlgoNomad", "BuildBro", "SaaSPhantom", "PromptPilot",
	"ShipShift", "LaunchLynx", "DebugDruid", "SyntaxSurfer", "DeployDino",
}

var handles = []string{
	"@byteWizard", "@cloudCrafter", "@pixelPulse", "@devSpecter", "@codeVoyager",
	"@stackSprite", "@algoNomad", "@buildBro", "@saasPhantom", "@promptPilot",
	"@shipShift", "@launchLynx", "@debugDruid", "@syntaxSurfer", "@deployDino",
}

func initGemini() error {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	geminiClient = client
	return nil
}

func generateWithGemini(ctx context.Context, categories []string, language string) (string, error) {
	// Create examples based on selected categories
	var examples []string
	if contains(categories, "saas") {
		examples = append(examples, saasExamples...)
	}
	if contains(categories, "ai") {
		examples = append(examples, aiExamples...)
	}
	if contains(categories, "growth") {
		examples = append(examples, growthExamples...)
	}

	if len(examples) == 0 {
		examples = append(examples, saasExamples...)
	}

	// Get language name for the prompt
	languageName, exists := languageNames[language]
	if !exists {
		languageName = "English"
		language = "en"
	}

	var languageInstruction string
	if language != "en" {
		languageInstruction = fmt.Sprintf(" Generate the post in %s. Make sure to adapt the tech bro culture and buzzwords to the target language while maintaining the satirical tone.", languageName)
	}

	systemPrompt := fmt.Sprintf(`You are a TechBro Post Generator. Generate a short, satirical tech bro social media post in the style of the examples provided. The post should be cringe-worthy, overly motivational, and include typical tech bro buzzwords and hashtags. Keep it authentic to the TechBro culture on social media.%s

Make sure to:
- Use short, punchy sentences
- Include relevant hashtags
- Sound overly confident and motivational
- Include buzzwords like "building", "scaling", "grinding", etc.
- Keep it under 280 characters
- Make it feel authentic but slightly exaggerated
%s

Examples:`, languageInstruction, 
		func() string {
			if language != "en" {
				return fmt.Sprintf("- Generate the content in %s", languageName)
			}
			return ""
		}())

	for _, example := range examples {
		systemPrompt += "\n\n" + example
	}

	userPrompt := "Generate a new TechBro post similar to these examples but completely original."

	combinedPrompt := fmt.Sprintf("System: %s\n\nUser: %s", systemPrompt, userPrompt)

	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				genai.NewPartFromText(combinedPrompt),
			},
		},
	}

	resp, err := geminiClient.Models.GenerateContent(ctx, "gemini-1.5-flash", contents, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate text with Gemini: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	textPart := resp.Candidates[0].Content.Parts[0]
	if textPart.Text == "" {
		return "", fmt.Errorf("no text content in Gemini response")
	}

	return strings.TrimSpace(textPart.Text), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func generatePost(c *gin.Context) {
	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate post text with Gemini
	postText, err := generateWithGemini(c.Request.Context(), req.Categories, req.Language)
	if err != nil {
		log.Printf("Error generating post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate post"})
		return
	}

	// Randomly select name, handle, and avatar
	nameIdx := rand.Intn(len(names))
	handleIdx := rand.Intn(len(handles))

	response := PostResponse{
		Name:     names[nameIdx],
		Handle:   handles[handleIdx],
		Avatar:   "/assets/profile.jpeg",
		PostText: postText,
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Initialize Gemini client
	if err := initGemini(); err != nil {
		log.Fatalf("Failed to initialize Gemini: %v", err)
	}

	// Set up Gin router
	router := gin.Default()

	// Enable CORS for development
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Serve static files
	router.Static("/assets", "./assets")
	router.StaticFile("/", "./index.html")

	// API routes
	router.POST("/api/generate", generatePost)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(router.Run(":" + port))
}
