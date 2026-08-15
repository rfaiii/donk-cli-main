import SwiftUI

struct HeaderView: View {
    var body: some View {
        HStack {
            Spacer()
            Text("//// DONK-CLI ////")
                .font(.custom("JetBrainsMonoNerdFontMono-Regular", size: 32))
                .foregroundColor(Theme.primary)
            Spacer()
        }
        .padding(.horizontal)
        .padding(.top, 10)
    }
}
