import SwiftUI

struct ContentView: View {
    @State private var mockStore = MockStore()

    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()
            
            ScrollView {
                VStack(spacing: 20) {
                    HeaderView()
                    
                    StatsView()
                    
                    TaskPickerView()
                    
                    ModelListView()
                    
                    BottomToolbarView()
                }
                .padding(.top, 10)
            }
        }
        .environment(mockStore)
        .preferredColorScheme(.dark)
    }
}

#Preview {
    ContentView()
}
