import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import PasteView from "../views/PasteView.vue";
import DocumentationView from "../views/DocumentationView.vue";
import MySharesView from "../views/MySharesView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", name: "home", component: HomeView },
    { path: "/p/:id", name: "paste", component: PasteView },
    { path: "/docs", name: "docs", component: DocumentationView },
    { path: "/my-shares", name: "myshares", component: MySharesView },
  ],
  scrollBehavior: () => ({ top: 0 }),
});

export default router;
