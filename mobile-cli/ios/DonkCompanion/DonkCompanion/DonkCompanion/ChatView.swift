import SwiftUI

struct ChatMessage: Identifiable {
    let id = UUID()
    let isUser: Bool
    let text: String
}

struct ChatView: View {
    @State private var messages: [ChatMessage] = [
        ChatMessage(isUser: true, text: "Can you help me parse this JSON file?"),
        ChatMessage(isUser: false, text: "Sure! Let's take a look at the structure of your JSON file. You can use the `jq` tool or write a small Go script to unmarshal it into a struct. What language are you using?"),
        ChatMessage(isUser: true, text: "I'm using Go.")
    ]
    
    var body: some View {
        ScrollViewReader { proxy in
            ScrollView {
                VStack(spacing: 16) {
                    ForEach(messages) { message in
                        HStack {
                            if message.isUser {
                                Spacer()
                                Text(message.text)
                                    .font(.custom("JetBrainsMonoNerdFontMono-Regular", size: 14))
                                    .padding()
                                    .background(Theme.surface)
                                    .foregroundColor(Theme.textPrimary)
                                    .cornerRadius(16)
                                    .overlay(
                                        RoundedRectangle(cornerRadius: 16)
                                            .stroke(Theme.border, lineWidth: 1)
                                    )
                                    .padding(.leading, 40)
                            } else {
                                Text(message.text)
                                    .font(.custom("JetBrainsMonoNerdFontMono-Regular", size: 14))
                                    .padding()
                                    .background(Color.clear)
                                    .foregroundColor(Theme.primary)
                                    .cornerRadius(16)
                                    .overlay(
                                        RoundedRectangle(cornerRadius: 16)
                                            .stroke(Theme.primary.opacity(0.3), lineWidth: 1)
                                    )
                                    .padding(.trailing, 40)
                                Spacer()
                            }
                        }
                        .id(message.id)
                    }
                }
                .padding(.horizontal)
                .padding(.vertical, 8)
            }
            .onAppear {
                if let last = messages.last {
                    proxy.scrollTo(last.id, anchor: .bottom)
                }
            }
        }
    }
}
