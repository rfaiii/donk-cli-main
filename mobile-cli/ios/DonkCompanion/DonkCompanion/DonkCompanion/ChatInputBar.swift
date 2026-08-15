import SwiftUI

struct ChatInputBar: View {
    @State private var text: String = ""
    
    var body: some View {
        HStack(spacing: 12) {
            // Upload button
            Button(action: {}) {
                Image(systemName: "plus")
                    .foregroundColor(Theme.primary)
                    .font(.system(size: 20))
            }
            
            // Text field
            TextField("Message DONK...", text: $text)
                .font(.custom("JetBrainsMonoNerdFontMono-Regular", size: 14))
                .foregroundColor(Theme.textPrimary)
                .padding(.vertical, 10)
                .padding(.horizontal, 12)
                .background(Theme.surface)
                .cornerRadius(20)
                .overlay(
                    RoundedRectangle(cornerRadius: 20)
                        .stroke(Theme.border, lineWidth: 1)
                )
            
            // Share button
            Button(action: {}) {
                Image(systemName: "square.and.arrow.up")
                    .foregroundColor(Theme.primary)
                    .font(.system(size: 20))
            }
        }
        .padding(.horizontal)
        .padding(.vertical, 8)
    }
}
