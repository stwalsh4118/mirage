import { ProjectDetail } from "@/components/project/project-detail"

interface ProjectPageProps {
  params: {
    id: string
  }
}

export default function ProjectPage({ params }: ProjectPageProps) {
  return (
    <div className="min-h-screen bg-background sandstorm-bg">
      <ProjectDetail projectId={params.id} />
    </div>
  )
}
