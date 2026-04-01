package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-bellringer/internal/store")
func(s *Server)handleSubscribe(w http.ResponseWriter,r *http.Request){var sub store.Subscription;json.NewDecoder(r.Body).Decode(&sub);if sub.Endpoint==""{writeError(w,400,"endpoint required");return};s.db.Subscribe(&sub);writeJSON(w,201,map[string]string{"status":"subscribed"})}
func(s *Server)handleUnsubscribe(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Unsubscribe(id);writeJSON(w,200,map[string]string{"status":"unsubscribed"})}
func(s *Server)handleListSubs(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListSubscriptions();if list==nil{list=[]store.Subscription{}};writeJSON(w,200,list)}
func(s *Server)handleSend(w http.ResponseWriter,r *http.Request){var n store.Notification;json.NewDecoder(r.Body).Decode(&n);if n.Title==""{writeError(w,400,"title required");return};result,_:=s.db.SendNotification(&n);writeJSON(w,200,result)}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
